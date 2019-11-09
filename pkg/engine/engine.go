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

package engine

import (
	"errors"

	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
)

const (
	componentName    = "engine"
	defaultFrequency = 1200
	maxFrequency     = 1450
	minFrequency     = 950
	maxSpeed         = 100
	minSpeed         = -100
)

type Engine struct {
	driver  *i2c.PCA9685Driver
	channel int
	logger  *logrus.Entry
	ready   bool
}

func NewEngine(driver *i2c.PCA9685Driver, channel int) *Engine {
	return &Engine{
		logger: logrus.WithFields(logrus.Fields{
			"group":     "car",
			"component": componentName,
		}),
		driver:  driver,
		channel: channel,
	}
}

// Init engine
func (e *Engine) Init() error {
	e.logger.Infof("init engine, channel: %d, frequency: %d", e.channel, defaultFrequency)

	if err := e.driver.SetPWM(e.channel, 0, defaultFrequency); err != nil {
		return err
	}

	e.ready = true

	return nil
}

// Move From -100 to 100
func (e *Engine) Move(speed int) error {
	if !e.ready {
		return errors.New("engine must be initiated before use")
	}

	frequency := uint16(gobot.ToScale(gobot.FromScale(float64(speed), minSpeed, maxSpeed), float64(minFrequency), float64(maxFrequency)))
	e.logger.Infof("move, frequency : %d, speed: %d, maxSpeed: %d, minSpeed: %d", frequency, speed, maxSpeed, minSpeed)

	return e.driver.SetPWM(e.channel, 0, frequency)
}
