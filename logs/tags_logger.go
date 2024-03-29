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
package logs

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger .
type Logger interface {
	SetLevel(level logrus.Level)
	IsLevelEnabled(level logrus.Level) bool
	Debug(args ...interface{})
	Info(args ...interface{})
	// Notice(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	// Critical(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	// Noticef(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	// Criticalf(format string, args ...interface{})
}

// TaggedLogger .
type TaggedLogger struct {
	logger    Logger
	tags      map[string]string
	tagsOrder []string
}

// SetTags .
func (c *TaggedLogger) setTags(tagValuePairs ...string) *TaggedLogger {
	for i := 1; i < len(tagValuePairs); i += 2 {
		c.setTag(tagValuePairs[i-1], tagValuePairs[i])
	}
	return c
}

// SetTag .
func (c *TaggedLogger) setTag(tag string, value interface{}) {
	tag = strings.Replace(tag, "%", "%%", -1)
	_, exists := c.tags[tag]
	c.tags[tag] = strings.Replace(fmt.Sprintf("%+v", value), "%", "%%", -1)
	if !exists {
		c.tagsOrder = append(c.tagsOrder, tag)
	}
}

// formatTags .
func (c *TaggedLogger) formatTags(prefix string) string {
	if len(c.tagsOrder) == 0 {
		return ""
	}
	// @TODO: Refactor to speedup
	var tagsString = prefix
	for inx, tag := range c.tagsOrder {
		if inx > 0 {
			tagsString += " "
		}
		tagsString += tag + `=[` + c.tags[tag] + `]`
	}
	return tagsString
}

// AddTags appends tags and create new instance.
func (c *TaggedLogger) AddTags(tagValuePairs ...string) *TaggedLogger {
	logger := TagLogger(c.logger)
	for tag, val := range c.tags {
		logger.setTag(tag, val)
	}
	logger.setTags(tagValuePairs...)
	return logger
}

// SetLevel .
func (c *TaggedLogger) SetLevel(level logrus.Level) {
	c.logger.SetLevel(level)
}

// IsEnabledFor .
func (c *TaggedLogger) IsLevelEnabled(level logrus.Level) bool {
	return c.logger.IsLevelEnabled(level)
}

// Debug .
func (c *TaggedLogger) Debug(args ...interface{}) {
	if c.IsLevelEnabled(logrus.DebugLevel) {
		c.logger.Debug(append(args, c.formatTags(""))...)
	}
}

// Info .
func (c *TaggedLogger) Info(args ...interface{}) {
	if c.IsLevelEnabled(logrus.InfoLevel) {
		c.logger.Info(append(args, c.formatTags(""))...)
	}
}

// Notice .
// func (c *TaggedLogger) Notice(args ...interface{}) {
// 	if c.IsEnabledFor(logrus.WarnLevel) {
// 		c.logger.Notice(append(args, c.formatTags(""))...)
// 	}
// }

// Warning .
func (c *TaggedLogger) Warning(args ...interface{}) {
	if c.IsLevelEnabled(logrus.WarnLevel) {
		c.logger.Warning(append(args, c.formatTags(""))...)
	}
}

// Error .
func (c *TaggedLogger) Error(args ...interface{}) {
	if c.IsLevelEnabled(logrus.ErrorLevel) {
		c.logger.Error(append(args, c.formatTags(""))...)
	}
}

// Critical .
// func (c *TaggedLogger) Critical(args ...interface{}) {
// 	if c.IsEnabledFor(logrus.FatalLevel) {
// 		c.logger.Critical(append(args, c.formatTags(""))...)
// 	}
// }

// Debugf .
func (c *TaggedLogger) Debugf(format string, args ...interface{}) {
	if c.IsLevelEnabled(logrus.DebugLevel) {
		c.logger.Debugf(format+c.formatTags(" "), args...)
	}
}

// Infof .
func (c *TaggedLogger) Infof(format string, args ...interface{}) {
	if c.IsLevelEnabled(logrus.InfoLevel) {
		c.logger.Infof(format+c.formatTags(" "), args...)
	}
}

// Noticef .
// func (c *TaggedLogger) Noticef(format string, args ...interface{}) {
// 	if c.IsEnabledFor(logrus.WarnLevel) {
// 		c.logger.Noticef(format+c.formatTags(" "), args...)
// 	}
// }

// Warningf .
func (c *TaggedLogger) Warningf(format string, args ...interface{}) {
	if c.IsLevelEnabled(logrus.WarnLevel) {
		c.logger.Warningf(format+c.formatTags(" "), args...)
	}
}

// Errorf .
func (c *TaggedLogger) Errorf(format string, args ...interface{}) {
	if c.IsLevelEnabled(logrus.ErrorLevel) {
		c.logger.Errorf(format+c.formatTags(" "), args...)
	}
}

// Criticalf .
// func (c *TaggedLogger) Criticalf(format string, args ...interface{}) {
// 	if c.IsEnabledFor(logrus.FatalLevel) {
// 		c.logger.Criticalf(format+c.formatTags(" "), args...)
// 	}
// }

// WithTags extends logger by tags.
func WithTags(logger Logger, tagValuePairs ...string) Logger {
	if tlog, ok := logger.(*TaggedLogger); ok {
		return tlog.AddTags(tagValuePairs...)
	}
	return TagLogger(logger, tagValuePairs...)
}

// TagLogger .
func TagLogger(logger Logger, tagValuePairs ...string) *TaggedLogger {
	return (&TaggedLogger{
		logger:    logger,
		tags:      make(map[string]string, 4),
		tagsOrder: make([]string, 0, 4),
	}).setTags(tagValuePairs...)
}
