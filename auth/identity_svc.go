package auth

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
)

//go:generate mockgen -source=identity_svc.go -package=auth -destination=identity_svc_mocks.go

type (
	// AttributeValue defines structure for attributes.
	AttributeValue struct {
		Value     string
		IsDefined bool
	}

	// IdentityService defines identity of the user/client.
	IdentityService interface {
		// ID returns id for group of users/clients
		MspID() (string, error)

		// CreatorID returns id for concrete user/client
		CreatorID() (string, error)

		// Cert returns user/client cert
		Cert() (*x509.Certificate, error)

		// CertID returns issue+subj of concrete user/client cert
		CertID() (string, error)

		// GetAttribute returns cert attribute
		GetAttribute(attrName string) (AttributeValue, error)
	}

	identityServiceImpl struct {
		stub      shim.ChaincodeStubInterface
		clientID  cid.ClientIdentity
		creatorID string
	}
)

func (svc *identityServiceImpl) init() error {
	if svc.clientID != nil {
		return nil
	}

	var err error

	svc.clientID, err = cid.New(svc.stub)
	if err != nil {
		return err
	}

	return nil
}

func (svc *identityServiceImpl) CertID() (string, error) {
	if err := svc.init(); err != nil {
		return "", err
	}

	return svc.clientID.GetID()
}

func (svc *identityServiceImpl) Cert() (*x509.Certificate, error) {
	if err := svc.init(); err != nil {
		return nil, err
	}

	return svc.clientID.GetX509Certificate()
}

func (svc *identityServiceImpl) CreatorID() (string, error) {
	if len(svc.creatorID) != 0 {
		return svc.creatorID, nil
	}

	bytes, err := svc.stub.GetCreator()
	if err != nil {
		return "", err
	}

	array := sha256.Sum256(bytes)
	svc.creatorID = hex.EncodeToString(array[:])

	return svc.creatorID, nil
}

func (svc *identityServiceImpl) MspID() (string, error) {
	if err := svc.init(); err != nil {
		return "", err
	}

	return svc.clientID.GetMSPID()
}

func (svc *identityServiceImpl) GetAttribute(attrName string) (out AttributeValue, err error) {
	var value string
	var find bool

	value, find, err = svc.clientID.GetAttributeValue(attrName)
	if err != nil {
		return
	}

	out = AttributeValue{value, find}

	return
}

// NewIdentityServiceImpl returns default implementation
func NewIdentityServiceImpl(stub shim.ChaincodeStubInterface) IdentityService {
	return &identityServiceImpl{
		stub: stub,
	}
}
