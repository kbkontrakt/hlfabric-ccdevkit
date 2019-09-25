package extstub

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type privateStubDecorator struct {
	shim.ChaincodeStubInterface
	collectionName string
}

func (s *privateStubDecorator) GetState(key string) ([]byte, error) {
	return s.ChaincodeStubInterface.GetPrivateData(s.collectionName, key)
}
func (s *privateStubDecorator) PutState(key string, value []byte) error {
	return s.ChaincodeStubInterface.PutPrivateData(s.collectionName, key, value)
}
func (s *privateStubDecorator) DelState(key string) error {
	return s.ChaincodeStubInterface.DelPrivateData(s.collectionName, key)
}
func (s *privateStubDecorator) SetStateValidationParameter(key string, ep []byte) error {
	return s.ChaincodeStubInterface.SetPrivateDataValidationParameter(s.collectionName, key, ep)
}
func (s *privateStubDecorator) GetStateValidationParameter(key string) ([]byte, error) {
	return s.ChaincodeStubInterface.GetPrivateDataValidationParameter(s.collectionName, key)
}
func (s *privateStubDecorator) GetStateByRange(startKey, endKey string) (shim.StateQueryIteratorInterface, error) {
	return s.ChaincodeStubInterface.GetPrivateDataByRange(s.collectionName, startKey, endKey)
}
func (s *privateStubDecorator) GetStateByRangeWithPagination(startKey, endKey string, pageSize int32,
	bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	panic("not supported")
}
func (s *privateStubDecorator) GetStateByPartialCompositeKey(objectType string, keys []string) (shim.StateQueryIteratorInterface, error) {
	return s.ChaincodeStubInterface.GetPrivateDataByPartialCompositeKey(s.collectionName, objectType, keys)
}
func (s *privateStubDecorator) GetStateByPartialCompositeKeyWithPagination(objectType string, keys []string,
	pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	panic("not supported")
}
func (s *privateStubDecorator) GetQueryResult(query string) (shim.StateQueryIteratorInterface, error) {
	return s.ChaincodeStubInterface.GetPrivateDataQueryResult(s.collectionName, query)
}
func (s *privateStubDecorator) GetQueryResultWithPagination(query string, pageSize int32,
	bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	panic("not supported")
}
func (s *privateStubDecorator) GetHistoryForKey(key string) (shim.HistoryQueryIteratorInterface, error) {
	return s.ChaincodeStubInterface.GetHistoryForKey(key)
}

// NewPrivateStubDecorator decorates stub for using private data collection as a source.
func NewPrivateStubDecorator(collectionName string, stub shim.ChaincodeStubInterface) shim.ChaincodeStubInterface {
	return &privateStubDecorator{
		ChaincodeStubInterface: stub,
		collectionName:         collectionName,
	}
}
