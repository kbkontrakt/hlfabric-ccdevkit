package auth

import "errors"

//go:generate mockgen -source=acl_svc.go -package=auth -destination=acl_svc_mocks.go

type (
	// AccessName .
	AccessName string

	// ControlList .
	ControlList struct {
		Attrs  map[string][]string `json:"attrs,omitempty"`
		MspID  map[string]bool     `json:"msp,omitempty"`
		CertID map[string]bool     `json:"cid,omitempty"`
	}

	// MatchAccessListFunc .
	MatchAccessListFunc func(name AccessName, ctrlList ControlList, idSvc IdentityService) error

	// ACL .
	ACL map[AccessName]ControlList
)

var (
	// ErrAccessNameNotFound .
	ErrAccessNameNotFound = errors.New("access name not found")
	// ErrAccessRestricted .
	ErrAccessRestricted = errors.New("access restricted")
)

// CombineMatchAccessAnd .
func CombineMatchAccessAnd(funcs ...MatchAccessListFunc) MatchAccessListFunc {
	return func(name AccessName, ctrlList ControlList, idSvc IdentityService) error {
		var err error
		for _, f := range funcs {
			err = f(name, ctrlList, idSvc)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// CombineMatchAccessOr .
func CombineMatchAccessOr(funcs ...MatchAccessListFunc) MatchAccessListFunc {
	return func(name AccessName, ctrlList ControlList, idSvc IdentityService) error {
		var err error
		for _, f := range funcs {
			err = f(name, ctrlList, idSvc)
			if err != ErrAccessRestricted {
				return err
			}
			if err == nil {
				return nil
			}
		}
		return ErrAccessRestricted
	}
}

// MatchAccessMspID .
func MatchAccessMspID() MatchAccessListFunc {
	return func(name AccessName, ctrlList ControlList, idSvc IdentityService) error {
		if ctrlList.MspID == nil {
			return ErrAccessRestricted
		}

		mspID, err := idSvc.MspID()
		if err != nil {
			return err
		}

		if ok, val := ctrlList.MspID[mspID]; !ok || !val {
			return ErrAccessRestricted
		}

		return nil
	}
}

// MatchAccessCertID .
func MatchAccessCertID() MatchAccessListFunc {
	return func(name AccessName, ctrlList ControlList, idSvc IdentityService) error {
		if ctrlList.CertID == nil {
			return ErrAccessRestricted
		}

		certID, err := idSvc.CertID()
		if err != nil {
			return err
		}

		if ok, val := ctrlList.CertID[certID]; !ok || !val {
			return ErrAccessRestricted
		}

		return nil
	}
}

// MatchAccessAttrs .
func MatchAccessAttrs() MatchAccessListFunc {
	return func(name AccessName, ctrlList ControlList, idSvc IdentityService) error {
		if ctrlList.Attrs == nil {
			return ErrAccessRestricted
		}

		var found bool

		for attrName, attrValues := range ctrlList.Attrs {
			idAttrVal, err := idSvc.GetAttribute(attrName)
			if err != nil {
				return err
			}
			if !idAttrVal.IsDefined {
				return ErrAccessRestricted
			}

			found = false
			for _, attrVal := range attrValues {
				if attrVal == idAttrVal.Value {
					found = true
					break
				}
			}
			if !found {
				return ErrAccessRestricted
			}
		}

		return nil
	}
}

// CheckAccess .
func (acl ACL) CheckAccess(name AccessName, idSvc IdentityService, matchFunc MatchAccessListFunc) error {
	list, found := acl[name]
	if !found {
		return ErrAccessNameNotFound
	}

	return matchFunc(name, list, idSvc)
}

type (
	// ACLService .
	ACLService interface {
		IsAllow(accessName string) (err error)
	}

	aclServiceImpl struct {
		aclRep    ACLRepository
		idSvc     IdentityService
		matchFunc MatchAccessListFunc
	}
)

// IsFuncAllow .
func (svc *aclServiceImpl) IsAllow(accessName string) error {
	acl, err := svc.aclRep.Get()
	if err != nil {
		return err
	}

	return acl.CheckAccess(AccessName(accessName), svc.idSvc, svc.matchFunc)
}

// NewACLServiceImpl creates default implementation
func NewACLServiceImpl(aclRep ACLRepository, idSvc IdentityService, matchFunc MatchAccessListFunc) ACLService {
	return &aclServiceImpl{aclRep, idSvc, matchFunc}
}