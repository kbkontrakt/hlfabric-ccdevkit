package extstub

import (
	"encoding/json"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	qr "github.com/hyperledger/fabric/protos/ledger/queryresult"
)

type (
	// PrivateHistoryStrategy .
	PrivateHistoryStrategy interface {
		// Append .
		Append(stub shim.ChaincodeStubInterface, collection, key string, value []byte, isDelete bool) error
		// GetIterator .
		GetIterator(stub shim.ChaincodeStubInterface, collection, key string) (shim.HistoryQueryIteratorInterface, error)
	}

	privateHistoryArrayAppendStrategy struct {
		keysPrefix string
		keysSuffix string
	}

	privateHistoryArrayAppendIterator struct {
		hist []qr.KeyModification
		inx  int
	}

	privateHistoryStubDecorator struct {
		shim.ChaincodeStubInterface
		history    PrivateHistoryStrategy
		collection string
	}

	keyValueHistory struct {
		TxID      string               `json:"i,omitempty"`
		Value     string               `json:"v,omitempty"`
		Timestamp *timestamp.Timestamp `json:"t,omitempty"`
		IsDelete  bool                 `json:"d,omitempty"`
	}
)

// Strategies
func (*privateHistoryArrayAppendStrategy) tryToFindPreviousActualItem(hist []keyValueHistory, item *keyValueHistory) {
	for inx := len(hist) - 1; inx >= 0; inx-- {
		if hist[inx].Value != "" {
			item.Value = hist[inx].Value
		}
	}
}

func (a *privateHistoryArrayAppendStrategy) Append(stub shim.ChaincodeStubInterface, collection, key string, value []byte, isDelete bool) error {
	timestamp, err := stub.GetTxTimestamp()
	if err != nil {
		return err
	}

	data, err := stub.GetPrivateData(collection, a.keysPrefix+key+a.keysSuffix)
	if err != nil {
		return err
	}

	// @TODO: add marshal strategy

	hist := []keyValueHistory{}
	if data != nil {
		err = json.Unmarshal(data, &hist)
		if err != nil {
			return err
		}
	}

	newHistItem := keyValueHistory{
		TxID:      stub.GetTxID(),
		Timestamp: timestamp,
		IsDelete:  isDelete,
		Value:     string(value),
	}
	if isDelete {
		// @TODO: Is need here?
		a.tryToFindPreviousActualItem(hist, &newHistItem)
	}

	hist = append([]keyValueHistory{newHistItem}, hist...)

	data, err = json.Marshal(hist)
	if err != nil {
		return err
	}

	return stub.PutPrivateData(collection, a.keysPrefix+key+a.keysSuffix, data)
}

func (a *privateHistoryArrayAppendStrategy) GetIterator(stub shim.ChaincodeStubInterface, collection, key string) (shim.HistoryQueryIteratorInterface, error) {
	data, err := stub.GetPrivateData(collection, a.keysPrefix+key+a.keysSuffix)
	if err != nil {
		return nil, err
	}

	// @TODO: add marshal strategy

	hist := []qr.KeyModification{}
	if data != nil {
		rawHist := []keyValueHistory{}
		err = json.Unmarshal(data, &rawHist)
		if err != nil {
			return nil, err
		}
		for _, item := range rawHist {
			hist = append(hist, qr.KeyModification{
				TxId:      item.TxID,
				Value:     []byte(item.Value),
				Timestamp: item.Timestamp,
				IsDelete:  item.IsDelete,
			})
		}
	}

	return &privateHistoryArrayAppendIterator{hist, len(hist)}, nil
}

func (i *privateHistoryArrayAppendIterator) HasNext() bool {
	return i.inx > 0
}
func (i *privateHistoryArrayAppendIterator) Next() (*qr.KeyModification, error) {
	if !i.HasNext() {
		return nil, nil
	}
	i.inx--
	return &i.hist[i.inx], nil
}
func (i *privateHistoryArrayAppendIterator) Close() error {
	return nil
}

// NewPrivateHistoryArrayAppendStrategy .
func NewPrivateHistoryArrayAppendStrategy(keysPrefix, keysSuffix string) PrivateHistoryStrategy {
	return &privateHistoryArrayAppendStrategy{
		keysPrefix: keysPrefix,
		keysSuffix: keysSuffix,
	}
}

// Stub

func (s *privateHistoryStubDecorator) PutPrivateData(collection, key string, value []byte) error {
	err := s.history.Append(s.ChaincodeStubInterface, s.collection, key, value, false)
	if err != nil {
		return err
	}
	return s.ChaincodeStubInterface.PutPrivateData(collection, key, value)
}
func (s *privateHistoryStubDecorator) DelPrivateData(collection, key string) error {
	err := s.history.Append(s.ChaincodeStubInterface, s.collection, key, nil, true)
	if err != nil {
		return err
	}
	return s.ChaincodeStubInterface.DelPrivateData(s.collection, key)
}
func (s *privateHistoryStubDecorator) GetHistoryForKey(key string) (shim.HistoryQueryIteratorInterface, error) {
	return s.history.GetIterator(s.ChaincodeStubInterface, s.collection, key)
}

// NewPrivateHistoryStubDecorator decorates stub for using private data collection with history request.
func NewPrivateHistoryStubDecorator(collectionName string, histStrategy PrivateHistoryStrategy, stub shim.ChaincodeStubInterface) shim.ChaincodeStubInterface {
	return &privateHistoryStubDecorator{
		collection: collectionName,
		history:    histStrategy,

		ChaincodeStubInterface: stub,
	}
}
