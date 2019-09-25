package debug

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var debugLog = shim.NewLogger("debug_stubs")

type debugGetStateByEmptyKey struct {
	shim.ChaincodeStubInterface
}

func (s *debugGetStateByEmptyKey) GetState(key string) ([]byte, error) {
	if key == "" {
		debugLog.Debugf("[DEBUG] Catch! GetState with empty key in [%s]\n", GetStacktrace(false))
	}
	return s.ChaincodeStubInterface.GetState(key)
}

func (s *debugGetStateByEmptyKey) GetPrivateData(collectionName, key string) ([]byte, error) {
	if key == "" {
		debugLog.Debugf("[DEBUG] Catch! GetPrivateData with empty key in [%s]\n", GetStacktrace(false))
	}
	return s.ChaincodeStubInterface.GetPrivateData(collectionName, key)
}

// NewDebugGetStateByEmptyKey notify places of passed empty keys with stacktrace.
func NewDebugGetStateByEmptyKey(stub shim.ChaincodeStubInterface) shim.ChaincodeStubInterface {
	return &debugGetStateByEmptyKey{
		ChaincodeStubInterface: stub,
	}
}
