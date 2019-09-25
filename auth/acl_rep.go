package auth

//go:generate mockgen -source=acl_rep.go -package=auth -destination=acl_rep_mocks.go

import (
	"errors"

	"kb-kontrakt.ru/hlfabric/ccdevkit/extstub"
)

type (
	// ACLRepository .
	ACLRepository interface {
		Get() (ACL, error)
		Save(acl ACL) error
	}

	aclRepositoryImpl struct {
		aclKey string
		stub   extstub.MarshalStub
	}

	aclDocument struct {
		ACList ACL `json:"acl"`
	}
)

const aclDefaultKeyName = "ACL"

func (rep *aclRepositoryImpl) Save(acl ACL) error {
	return rep.stub.WriteState(rep.aclKey, aclDocument{acl})
}

func (rep *aclRepositoryImpl) Get() (ACL, error) {
	var acl aclDocument

	err := rep.stub.ReadState(rep.aclKey, &acl)
	if err == extstub.ErrNotFound {
		return nil, errors.New("ACL not found")
	}
	if err != nil {
		return nil, err
	}

	return acl.ACList, nil
}

// NewACLRepositoryImpl creates default acl implementation
func NewACLRepositoryImpl(stub extstub.MarshalStub) ACLRepository {
	return &aclRepositoryImpl{
		aclDefaultKeyName,
		stub,
	}
}
