package chrono

//go:generate mockgen -source=time_svc.go -package=chrono -destination=time_svc_mocks.go

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type (
	// TimeService .
	TimeService interface {
		NowDateTime() (DateTime, error)
	}

	timeServiceImpl struct {
		stub shim.ChaincodeStubInterface
	}
)

func (ts *timeServiceImpl) NowDateTime() (dt DateTime, err error) {
	var txTime *timestamp.Timestamp

	txTime, err = ts.stub.GetTxTimestamp()
	if err != nil {
		return
	}

	return DateTimeFromTimestamp(txTime), nil
}

// NewTimeServiceImpl creates default implementation
func NewTimeServiceImpl(stub shim.ChaincodeStubInterface) TimeService {
	return &timeServiceImpl{
		stub,
	}
}
