package debug

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Invoke processes debug methods
func Invoke(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	const lockKey = "debuglock_"
	const batchEnumSize = 2048

	logger := shim.NewLogger("debugtools")

	lockData, err := stub.GetState(lockKey)
	if err != nil {
		return nil, err
	}
	if lockData != nil {
		return nil, errors.New("unsupported function")
	}

	tmap, err := stub.GetTransient()
	if tmap != nil && len(tmap["args"]) > 0 {
		argsInfs := []interface{}{}
		err = json.Unmarshal(tmap["args"], &argsInfs)
		if err != nil {
			return nil, err
		}
		args = []string{}
		for _, arg := range argsInfs {
			switch val := arg.(type) {
			case string:
				args = append(args, val)
			default:
				bytes, err := json.Marshal(arg)
				if err != nil {
					return nil, err
				}
				args = append(args, string(bytes))
			}
		}
	}

	var isForPrivateData = false
	var delKeysStartSliceIndex = 2

	switch args[0] {
	case "DelPrivState":
		isForPrivateData = true
		delKeysStartSliceIndex = 3
		fallthrough

	case "DelState":
		if len(args) < delKeysStartSliceIndex+1 {
			return nil, errors.New("not enough arguments")
		}
		var err error
		var keys []string

		if args[1] == "query" {
			var iterator shim.StateQueryIteratorInterface
			if isForPrivateData {
				// iterator, err = stub.GetPrivateDataQueryResult(args[2], args[3])
				panic("it is not allowed to query and update in the same transaction")
			} else {
				iterator, err = stub.GetQueryResult(args[2])
			}
			if err != nil {
				return nil, err
			}
			keys, err = enumKeysFromIterator(iterator, batchEnumSize)
		} else {
			keys, err = enumKeysFromJSONArgs(args[delKeysStartSliceIndex:], batchEnumSize)
		}
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			if isForPrivateData {
				err = stub.DelPrivateData(args[2], key)
			} else {
				err = stub.DelState(key)
			}

			if err != nil {
				logger.Errorf("failed to delete state of [%+v]: %+v", key, err)
			}
		}

	case "Lock":
		err = stub.PutState(lockKey, []byte("1"))
		break

	default:
		return nil, errors.New("unknown debug method")
	}

	return nil, nil
}
