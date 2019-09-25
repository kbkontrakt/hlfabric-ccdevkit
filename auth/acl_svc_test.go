package auth

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAclServiceSpec(t *testing.T) {
	Convey("Given pre-filled acl repository with allowed method1 for org1", t, func() {
		ctrl := gomock.NewController(t)

		aclRep := NewMockACLRepository(ctrl)
		aclRep.EXPECT().Get().
			Return(ACL{"method1": ControlList{
				MspID:  map[string]bool{"org1": true, "org2": false},
				CertID: map[string]bool{"cert1": true, "cert2": false},
				Attrs:  map[string][]string{"role": {"role1", "role2"}},
			}}, nil).Times(1)

		idSvc := NewMockIdentityService(ctrl)

		svc := NewACLServiceImpl(aclRep, idSvc, MatchAccessMspID())

		// msp

		Convey("When check access for method1 by msp org1", func() {
			idSvc.EXPECT().MspID().Return("org1", nil).Times(1)

			err := svc.IsAllow("method1")

			Convey("Then it should allow", func() {
				So(err, ShouldBeNil)
				ctrl.Finish()
			})
		})

		Convey("When check access for method1 by msp org2", func() {
			idSvc.EXPECT().MspID().Return("org2", nil).Times(1)

			err := svc.IsAllow("method1")

			Convey("Then it should restrict access", func() {
				So(err, ShouldBeError, ErrAccessRestricted)
				ctrl.Finish()
			})
		})

		Convey("When check access for method1 by msp org3", func() {
			idSvc.EXPECT().MspID().Return("org3", nil).Times(1)

			err := svc.IsAllow("method1")

			Convey("Then it should restrict access", func() {
				So(err, ShouldBeError, ErrAccessRestricted)
				ctrl.Finish()
			})
		})

		// cert

		svc = NewACLServiceImpl(aclRep, idSvc, MatchAccessCertID())

		Convey("When check access for method1 by cert1", func() {
			idSvc.EXPECT().CertID().Return("cert1", nil).Times(1)

			err := svc.IsAllow("method1")

			Convey("Then it should allow", func() {
				So(err, ShouldBeNil)
				ctrl.Finish()
			})
		})

		Convey("When check access for method1 by cert2", func() {
			idSvc.EXPECT().CertID().Return("cert2", nil).Times(1)

			err := svc.IsAllow("method1")

			Convey("Then it should restrict access", func() {
				So(err, ShouldBeError, ErrAccessRestricted)
				ctrl.Finish()
			})
		})

		Convey("When check access for method1 by cert3", func() {
			idSvc.EXPECT().CertID().Return("cert3", nil).Times(1)

			err := svc.IsAllow("method1")

			Convey("Then it should restrict access", func() {
				So(err, ShouldBeError, ErrAccessRestricted)
				ctrl.Finish()
			})
		})

		// attrs

		svc = NewACLServiceImpl(aclRep, idSvc, MatchAccessAttrs())

		Convey("When check access for method1 by role2", func() {
			idSvc.EXPECT().GetAttribute("role").Return(AttributeValue{"role2", true}, nil).Times(1)

			err := svc.IsAllow("method1")

			Convey("Then it should allow", func() {
				So(err, ShouldBeNil)
				ctrl.Finish()
			})
		})

		Convey("When check access for method1 without role attribute", func() {
			idSvc.EXPECT().GetAttribute("role").Return(AttributeValue{"", false}, nil).Times(1)

			err := svc.IsAllow("method1")

			Convey("Then it should restrict access", func() {
				So(err, ShouldBeError, ErrAccessRestricted)
				ctrl.Finish()
			})
		})

		// combine and

		svc = NewACLServiceImpl(aclRep, idSvc, CombineMatchAccessAnd(MatchAccessMspID(), MatchAccessAttrs()))

		Convey("When check access for method1 by msp org1 with role2", func() {
			idSvc.EXPECT().MspID().Return("org1", nil).Times(1)
			idSvc.EXPECT().GetAttribute("role").Return(AttributeValue{"role2", true}, nil).Times(1)

			err := svc.IsAllow("method1")

			Convey("Then it should allow", func() {
				So(err, ShouldBeNil)
				ctrl.Finish()
			})
		})

		Convey("When check access for method1 by msp org1 with role3", func() {
			idSvc.EXPECT().MspID().Return("org1", nil).Times(1)
			idSvc.EXPECT().GetAttribute("role").Return(AttributeValue{"role3", true}, nil).Times(1)

			err := svc.IsAllow("method1")

			Convey("Then it should restrict access", func() {
				So(err, ShouldBeError, ErrAccessRestricted)
				ctrl.Finish()
			})
		})
	})
}

func TestAclServiceInternalErrorSpec(t *testing.T) {
	Convey("Given pre-filled acl repository without any controllist for method1", t, func() {
		ctrl := gomock.NewController(t)

		aclRep := NewMockACLRepository(ctrl)
		aclRep.EXPECT().Get().
			Return(ACL{"method1": ControlList{}}, nil).Times(1)

		idSvc := NewMockIdentityService(ctrl)
		svc := NewACLServiceImpl(aclRep, idSvc, MatchAccessMspID())

		Convey("When check access for method1 by org1", func() {
			err := svc.IsAllow("method1")

			Convey("Then it should deny", func() {
				So(err, ShouldBeError, ErrAccessRestricted)
				ctrl.Finish()
			})
		})
	})

	Convey("Given acl repository that returns some error on request", t, func() {
		ctrl := gomock.NewController(t)

		aclRep := NewMockACLRepository(ctrl)
		aclRep.EXPECT().Get().
			Return(nil, errors.New("some error")).Times(1)

		idSvc := NewMockIdentityService(ctrl)
		svc := NewACLServiceImpl(aclRep, idSvc, MatchAccessMspID())

		Convey("When check access for method1 by org1", func() {
			err := svc.IsAllow("method1")

			Convey("Then it should passthrough that error", func() {
				So(err, ShouldBeError, "some error")
				ctrl.Finish()
			})
		})
	})

	Convey("Given identity service that returns some error on request mspID", t, func() {
		ctrl := gomock.NewController(t)

		aclRep := NewMockACLRepository(ctrl)
		aclRep.EXPECT().Get().
			Return(ACL{"method1": ControlList{
				MspID:  map[string]bool{"org1": true, "org2": false},
				CertID: map[string]bool{"cert1": true, "cert2": false},
				Attrs:  map[string][]string{"role": {"role1", "role2"}},
			}}, nil).Times(1)

		idSvc := NewMockIdentityService(ctrl)
		svc := NewACLServiceImpl(aclRep, idSvc, MatchAccessMspID())

		Convey("When check access for method1", func() {
			idSvc.EXPECT().MspID().Return("", errors.New("some error"))

			err := svc.IsAllow("method1")

			Convey("Then it should passthrough that error", func() {
				So(err, ShouldBeError, "some error")
				ctrl.Finish()
			})
		})
	})
}
