package logs

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// dummyLogger .
type dummyLogger struct{}

// SetLevel .
func (c *dummyLogger) SetLevel(level shim.LoggingLevel) {}

// IsEnabledFor .
func (c *dummyLogger) IsEnabledFor(level shim.LoggingLevel) bool {
	return true
}

// Debug .
func (c *dummyLogger) Debug(args ...interface{}) {}

// Info .
func (c *dummyLogger) Info(args ...interface{}) {}

// Notice .
func (c *dummyLogger) Notice(args ...interface{}) {}

// Warning .
func (c *dummyLogger) Warning(args ...interface{}) {}

// Error .
func (c *dummyLogger) Error(args ...interface{}) {}

// Critical .
func (c *dummyLogger) Critical(args ...interface{}) {}

// Debugf .
func (c *dummyLogger) Debugf(format string, args ...interface{}) {}

// Infof .
func (c *dummyLogger) Infof(format string, args ...interface{}) {}

// Noticef .
func (c *dummyLogger) Noticef(format string, args ...interface{}) {}

// Warningf .
func (c *dummyLogger) Warningf(format string, args ...interface{}) {}

// Errorf .
func (c *dummyLogger) Errorf(format string, args ...interface{}) {}

// Criticalf .
func (c *dummyLogger) Criticalf(format string, args ...interface{}) {}

// DummyLogger .
func DummyLogger() Logger {
	return &dummyLogger{}
}
