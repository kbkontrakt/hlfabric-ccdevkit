/*
 *  Copyright 2017 - 2019 KB Kontrakt LLC - All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */
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
