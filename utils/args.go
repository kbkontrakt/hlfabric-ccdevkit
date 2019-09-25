package utils

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// GetArgsFromTransientMap returns function name and args from transientMap
func GetArgsFromTransientMap(stub shim.ChaincodeStubInterface) (fn string, args []string, err error) {
	args = []string{}

	tmap, err := stub.GetTransient()
	if err != nil {
		return
	}

	data, exists := tmap["Args"]
	if !exists {
		return
	}

	rawArgsItems := []interface{}{}

	err = json.Unmarshal(data, &rawArgsItems)
	if err != nil {
		return
	}

	if len(rawArgsItems) == 0 {
		return
	}

	bytes := []byte{}
	for inx, arg := range rawArgsItems {
		argStr := ""
		if arg != nil {
			if val, ok := arg.(string); ok {
				argStr = val
			} else {
				bytes, err = json.Marshal(arg)
				if err != nil {
					return
				}
				argStr = string(bytes)
			}
		}

		if inx == 0 {
			fn = argStr
		} else {
			args = append(args, argStr)
		}
	}

	return
}

// GetFnArgsOrFromTransientMap returns function name and args from args else transientMap
func GetFnArgsOrFromTransientMap(stub shim.ChaincodeStubInterface) (fn string, args []string, err error) {
	fn, args = stub.GetFunctionAndParameters()

	if (len(fn) == 0 || fn == "*") && len(args) == 0 {
		fn, args, err = GetArgsFromTransientMap(stub)
	}

	return
}
