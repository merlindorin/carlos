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

package direction

import (
	"errors"

	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
)

const (
	componentName    = "direction"
	defaultFrequency = 1200
	maxFrequency     = 1450
	minFrequency     = 950
	maxRadius        = 90
	minRadius        = -90
)

type Direction struct {
	driver  *i2c.PCA9685Driver
	logger  *logrus.Entry
	channel int
	ready   bool
}

func NewDirection(driver *i2c.PCA9685Driver, channel int) *Direction {
	return &Direction{
		logger: logrus.WithFields(logrus.Fields{
			"group":     "car",
			"component": componentName,
		}),
		channel: channel,
		driver:  driver,
	}
}

// Init direction
func (e *Direction) Init() error {
	e.logger.Infof("init direction, channel: %d, frequency: %d", e.channel, defaultFrequency)
	if err := e.driver.SetPWM(e.channel, 0, defaultFrequency); err != nil {
		return err
	}

	e.ready = true

	return nil
}

// Move From -90 to 90
func (e *Direction) Turn(radius int) error {
	if !e.ready {
		return errors.New("direction must be initiated before use")
	}

	frequency := uint16(gobot.ToScale(gobot.FromScale(float64(radius), minRadius, maxRadius), float64(minFrequency), float64(maxFrequency)))
	e.logger.Infof("turn, frequency: %d, radius: %d, maxRadius: %d, minRadius: %d", frequency, radius, maxRadius, minRadius)

	return e.driver.SetPWM(e.channel, 0, frequency)
}
