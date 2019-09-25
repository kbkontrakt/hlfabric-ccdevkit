package tests

//go:generate mockgen -source=../vendor/github.com/hyperledger/fabric/core/chaincode/shim/interfaces.go -package=tests -destination=chaincode_mocks.go
//go:generate sed "10 i    \"github.com/hyperledger/fabric/core/chaincode/shim\"" -i chaincode_mocks.go
//go:generate sed -E -e "s/ ChaincodeStubInterface/ shim.ChaincodeStubInterface/g" -e "s/\\(StateQueryIteratorInterface/ (shim.StateQueryIteratorInterface/g" -e "s/\\(HistoryQueryIteratorInterface/(shim.HistoryQueryIteratorInterface/g" -i chaincode_mocks.go
//go:generate go fmt chaincode_mocks.go
