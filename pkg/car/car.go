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

package car

import (
	"github.com/iam-merlin/carlos/pkg/direction"
	"github.com/iam-merlin/carlos/pkg/engine"
	"time"

	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	startingDelay = 1 * time.Second
)

type Car struct {
	PCA9685Driver *i2c.PCA9685Driver
	Engine        *engine.Engine
	Direction     *direction.Direction
	Power         bool
	logger        *logrus.Entry
}

func NewCar() (*Car, error) {
	adaptor := raspi.NewAdaptor()
	driver := i2c.NewPCA9685Driver(adaptor)

	car := &Car{
		PCA9685Driver: driver,
		Direction:     direction.NewDirection(driver, 1),
		Engine:        engine.NewEngine(driver, 0),
		logger: logrus.WithFields(logrus.Fields{
			"group":     "car",
			"compoment": "main",
		}),
	}

	return car, nil
}

func (c *Car) init() error {
	if err := c.Direction.Init(); err != nil {
		return err
	}

	if err := c.Engine.Init(); err != nil {
		return err
	}

	return nil
}

// Power On
func (c *Car) PowerOn() error {
	c.logger.Info("Power ON")

	if c.Power {
		return nil
	}

	c.logger.Infof("Start i2c PCA9685 driver")
	if err := c.PCA9685Driver.Start(); err != nil {
		return err
	}

	if err := c.init(); err != nil {
		return err
	}

	c.logger.Infof("Wait %s for power ON", startingDelay)
	time.Sleep(startingDelay)
	c.Power = true

	return nil
}

// Power Off
func (c *Car) PowerOff() error {
	c.logger.Info("Power Off")

	if !c.Power {
		return nil
	}

	if err := c.init(); err != nil {
		return err
	}

	c.logger.Infof("Stop i2c PCA9685 PCA9685Driver")
	if err := c.PCA9685Driver.Halt(); err != nil {
		return err
	}

	c.logger.Infof("Wait %s for power Off", startingDelay)
	time.Sleep(startingDelay)

	c.Power = false

	return nil
}
