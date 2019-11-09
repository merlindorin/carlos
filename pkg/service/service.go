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

package service

import (
	"context"
	"errors"
	"github.com/iam-merlin/carlos"
	"github.com/iam-merlin/carlos/pkg/car"
	"github.com/iam-merlin/carlos/pkg/log"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/health"
)

//CarServiceImpl is a implementation of RepositoryService Grpc Service.
type CarServiceImpl struct {
	Health      *health.Server
	logChannel  *log.ChannelLogger
	serviceName string
	car         *car.Car
}

//NewCarServiceImpl returns the pointer to the implementation.
func NewCarServiceImpl(serviceName string, healthServer *health.Server) (*CarServiceImpl, error) {
	car, err := car.NewCar()
	if err != nil {
		return nil, err
	}

	// register log channel
	logChannel := log.NewChannelLogger()
	logrus.AddHook(logChannel)

	return &CarServiceImpl{
		car:         car,
		Health:      healthServer,
		serviceName: serviceName,
		logChannel:  logChannel,
	}, nil
}

func (serviceImpl *CarServiceImpl) End() error {
	return serviceImpl.car.PowerOff()
}

//Add function implementation of gRPC Service.
func (serviceImpl *CarServiceImpl) Log(empty *main.Empty, stream main.LogService_LogServer) error {
	serviceImpl.logChannel.IsSubscribing(true)
	defer serviceImpl.logChannel.IsSubscribing(false)

	for {
		select {
		case l := <-serviceImpl.logChannel.Channel:
			if err := stream.Send(&main.Log{Message: l.Message, Level: int64(l.Level)}); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}

// Power
func (serviceImpl *CarServiceImpl) Power(c context.Context, p *main.Power) (*main.Empty, error) {
	if p.Power {
		return &main.Empty{}, serviceImpl.car.PowerOn()
	}
	return &main.Empty{}, serviceImpl.car.PowerOff()
}

// Emergency
func (serviceImpl *CarServiceImpl) Emergency(context.Context, *main.Emergency) (*main.Empty, error) {
	if !serviceImpl.car.Power {
		return &main.Empty{}, errors.New("car must be power on before emergency use")
	}

	return &main.Empty{}, serviceImpl.car.PowerOff()
}

// Brake
func (serviceImpl *CarServiceImpl) Brake(c context.Context, b *main.Brake) (*main.Empty, error) {
	if !serviceImpl.car.Power {
		return &main.Empty{}, errors.New("car must be power on before brake use")
	}

	return &main.Empty{}, serviceImpl.car.Direction.Turn(int(b.Radius))
}

// Move
func (serviceImpl *CarServiceImpl) Move(c context.Context, m *main.Move) (*main.Empty, error) {
	if !serviceImpl.car.Power {
		return &main.Empty{}, errors.New("car must be power on before move use")
	}

	return &main.Empty{}, serviceImpl.car.Engine.Move(int(m.Speed))
}
