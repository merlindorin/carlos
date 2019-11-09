/*
* Copyright The Carlos Authors.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package log

import (
	"github.com/sirupsen/logrus"
)

type ChannelLogger struct {
	Channel        chan logrus.Entry
	withSubscriber bool
	levels         []logrus.Level
}

func NewChannelLogger() *ChannelLogger {
	return &ChannelLogger{
		Channel:        make(chan logrus.Entry, 50),
		withSubscriber: false,
		levels:         logrus.AllLevels,
	}
}

// A hook to be fired when logging
func (l *ChannelLogger) Fire(entry *logrus.Entry) error {
	// We push entries if we subscribe
	if l.withSubscriber {
		l.Channel <- *entry
	}

	return nil
}

// Level accepted
func (l *ChannelLogger) Levels() []logrus.Level {
	return l.levels
}

// Set acceptable levels
func (l *ChannelLogger) SetLevels(levels []logrus.Level) {
	l.levels = levels
}

// Set if we subscribe or not
func (l *ChannelLogger) IsSubscribing(subscribe bool) {
	l.withSubscriber = subscribe
}
