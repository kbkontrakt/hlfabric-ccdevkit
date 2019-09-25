package tests

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	gomock "github.com/golang/mock/gomock"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	queryresult "github.com/hyperledger/fabric/protos/ledger/queryresult"
)

// NewSMC .
func NewSMC(t *testing.T, ccname string) (*MockChaincodeStubInterface, *shim.MockStub, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mock := shim.NewMockStub(ccname, nil)
	mock.MockTransactionStart("TX1")

	return NewMockChaincodeStubInterface(ctrl), mock, ctrl
}

// GomockRexp .
func GomockRexp(rexp string) gomock.Matcher {
	return gomockRexpMatcher{regexp.MustCompile(rexp)}
}

type gomockRexpMatcher struct {
	rexp *regexp.Regexp
}

func (e gomockRexpMatcher) Matches(x interface{}) bool {
	return e.rexp.MatchString(fmt.Sprintf("%+v", x))
}

func (e gomockRexpMatcher) String() string {
	return fmt.Sprintf("is %v regexp matched", e.rexp)
}

func getSerializedValueObject(val interface{}) []byte {
	var err error
	var data []byte

	switch val := val.(type) {
	case []byte:
		data = val
	default:
		data, err = json.Marshal(&val)
		if err != nil {
			panic(err)
		}
	}

	return data
}

func jsonPreFillByKeyValPairs(callback func(string, []byte), keyValPairs ...interface{}) {
	for inx, l := 0, len(keyValPairs); inx < l; inx += 2 {
		keyName := keyValPairs[inx].(string)
		data := getSerializedValueObject(keyValPairs[inx+1])
		callback(keyName, data)
	}
}

// JSONFillMockState .
func JSONFillMockState(mock *shim.MockStub, keyValPairs ...interface{}) {
	jsonPreFillByKeyValPairs(func(key string, data []byte) {
		mock.PutState(key, data)
	}, keyValPairs...)
}

// JSONFillMockPrvState .
func JSONFillMockPrvState(mock *shim.MockStub, collectionName string, keyValPairs ...interface{}) {
	jsonPreFillByKeyValPairs(func(key string, data []byte) {
		mock.PutPrivateData(collectionName, key, data)
	}, keyValPairs...)
}

// WrapShimMockGetPrivState .
func WrapShimMockGetPrivState(mock *shim.MockStub) interface{} {
	return func(collection, key string) ([]byte, error) {
		return mock.GetPrivateData(collection, key)
	}
}

// WrapShimMockPutPrivState .
func WrapShimMockPutPrivState(mock *shim.MockStub) interface{} {
	return func(collection, key string, data []byte) error {
		return mock.PutPrivateData(collection, key, data)
	}
}

// WrapShimMockDelPrivState .
func WrapShimMockDelPrivState(mock *shim.MockStub) interface{} {
	return func(collection, key string, data []byte) error {
		return mock.DelPrivateData(collection, key)
	}
}

// WrapShimMockGetState .
func WrapShimMockGetState(mock *shim.MockStub) interface{} {
	return func(key string) ([]byte, error) {
		return mock.GetState(key)
	}
}

// WrapShimMockPutState .
func WrapShimMockPutState(mock *shim.MockStub) interface{} {
	return func(key string, data []byte) error {
		return mock.PutState(key, data)
	}
}

// WrapShimMockDelState .
func WrapShimMockDelState(mock *shim.MockStub) interface{} {
	return func(key string, data []byte) error {
		return mock.DelState(key)
	}
}

// WrapShimMockTxStamp .
func WrapShimMockTxStamp(mock *shim.MockStub) interface{} {
	return func() (*timestamp.Timestamp, error) {
		return mock.GetTxTimestamp()
	}
}

// WrapShimMockTxID .
func WrapShimMockTxID(mock *shim.MockStub) interface{} {
	return func() string {
		return mock.GetTxID()
	}
}

// MakeJSONStateIteratorFuncs .
func MakeJSONStateIteratorFuncs(namespace string, keyValPairs ...interface{}) (interface{}, interface{}) {
	inx, l := 0, len(keyValPairs)

	return func() (*queryresult.KV, error) {
			inx += 2
			if inx <= l {
				keyName := keyValPairs[inx-2].(string)
				return &queryresult.KV{
					Namespace: namespace,
					Key:       keyName,
					Value:     getSerializedValueObject(keyValPairs[inx-1]),
				}, nil
			}
			return nil, nil
		}, func() bool {
			if inx >= l {
				return false
			}
			return true
		}
}

// MakeJSONStateIterator .
func MakeJSONStateIterator(ctrl *gomock.Controller, namespace string,
	keyValPairs ...interface{}) *MockStateQueryIteratorInterface {

	nextTimes := 1
	if len(keyValPairs) == 0 {
		nextTimes = 0
	}

	iter, hasNext := MakeJSONStateIteratorFuncs(namespace, keyValPairs...)

	queryMock := NewMockStateQueryIteratorInterface(ctrl)
	queryMock.EXPECT().Next().DoAndReturn(iter).MinTimes(nextTimes)
	queryMock.EXPECT().HasNext().DoAndReturn(hasNext).MinTimes(1)
	queryMock.EXPECT().Close().MinTimes(1)

	return queryMock
}