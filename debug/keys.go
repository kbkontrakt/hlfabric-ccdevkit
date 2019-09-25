package debug

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type structBasedKey struct {
	IDsmall *string `json:"id"`
	IDbig   *string `json:"ID"`
}

func enumKeysFromJSONArgs(args []string, limit int) ([]string, error) {
	var err error
	keys := make([]string, 0, limit)

	for _, arg := range args {
		if len(arg) == 0 {
			continue
		}
		if arg[0] == '"' {
			key := ""
			err = json.Unmarshal([]byte(arg), &key)
			if err != nil {
				return nil, err
			}
			keys = append(keys, key)
			continue
		}

		structs := []structBasedKey{}
		if arg[0] == '[' {
			err = json.Unmarshal([]byte(arg), &structs)
		} else {
			keyStruct := structBasedKey{}
			err = json.Unmarshal([]byte(arg), &keyStruct)
			structs = append(structs, keyStruct)
		}

		for i, l := 0, len(structs); i < l && limit >= 0; i, limit = i+1, limit-1 {
			if structs[i].IDsmall != nil {
				keys = append(keys, *structs[i].IDsmall)
			} else if structs[i].IDbig != nil {
				keys = append(keys, *structs[i].IDbig)
			}
		}
	}

	return keys, nil
}

func enumKeysFromIterator(iter shim.StateQueryIteratorInterface, limit int) ([]string, error) {
	keys := make([]string, 0, limit)

	for ; iter.HasNext() && limit >= 0; limit-- {
		pair, err := iter.Next()

		if err != nil {
			return nil, err
		}
		keys = append(keys, pair.GetKey())
	}

	return keys, nil
}
