package utils

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSyncEventSpec(t *testing.T) {
	Convey("SyncEvent Spec", t, func(c C) {
		c.Convey("Given empty event", func(c C) {
			event := NewSyncEvent()
			c.Convey("When emit event", func(c C) {
				err := event.Emit("arg")
				c.Convey("It should be done without error", func(c C) {
					So(err, ShouldBeNil)
				})
			})
		})

		c.Convey("Given event with one listener", func(c C) {
			event := NewSyncEvent()
			var passedArg interface{}
			event.On(func(args ...interface{}) error {
				if len(args) == 1 {
					passedArg = args[0]
				}
				return nil
			})
			c.Convey("When emit event", func(c C) {
				err := event.Emit("arg1")
				c.Convey("It should receive value and done without error", func(c C) {
					So(err, ShouldBeNil)
					So(passedArg, ShouldEqual, "arg1")
				})
			})
		})

		c.Convey("Given event with wrapped to single json one listener", func(c C) {
			event := NewSyncEvent()
			var passedArg interface{}
			event.On(WrapOnSingleArg(MarshalFuncJSON, func(args ...interface{}) error {
				if len(args) == 1 {
					passedArg = args[0]
				}
				return nil
			}))
			c.Convey("When emit event with object argument", func(c C) {
				err := event.Emit(struct{ Value string }{"test"})
				c.Convey("It should receive value and done without error", func(c C) {
					So(err, ShouldBeNil)
					So(passedArg, ShouldResemble, []byte(`{"Value":"test"}`))
				})
			})
			c.Convey("When emit 2 events", func(c C) {
				err := event.Emit(struct{ Value string }{"test"}, "other")
				c.Convey("It should return error", func(c C) {
					So(err, ShouldBeError, "failed listen event: not single argument was passed")
				})
			})
		})

		c.Convey("Given event with one listener that returns error", func(c C) {
			event := NewSyncEvent()
			event.On(func(args ...interface{}) error {
				return errors.New("error")
			})
			c.Convey("When emit event", func(c C) {
				err := event.Emit("arg1")
				c.Convey("It should returns error from listener", func(c C) {
					So(err, ShouldBeError, "error")
				})
			})
		})
	})
}
