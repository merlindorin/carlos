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

package main

import (
	"context"
	"fmt"
	"github.com/iam-merlin/carlos/pkg/service"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	openCencusGrpc "github.com/akhenakh/ocgrpc_propagation"
	nested "github.com/antonfisher/nested-logrus-formatter"
	grpc2 "github.com/iam-merlin/carlos/internal/grpc"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthPB "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
)

const (
	version           = "dev"
	serviceName       = "car"
	defaultHost       = "0.0.0.0"
	defaultGrpcPort   = 10000
	defaultHealthPort = 9000
	defaultMetricPort = 8000
	httpReadTimeout   = 10 * time.Second
	httpWriteTimeout  = 10 * time.Second
)

var (
	grpcServer       *grpc.Server
	grpcHealthServer *grpc.Server
	httpServer       *http.Server

	kaep = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}

	kasp = keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
		MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
		MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "cartel"
	app.Description = "Cartel connect car devices into a group of grpc services"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "tls",
		},
		cli.StringFlag{
			Name: "cert-file",
		},
		cli.StringFlag{
			Name: "key-file",
		},
		cli.IntFlag{
			Name:  "http-metric-port",
			Value: defaultMetricPort,
		},
		cli.IntFlag{
			Name:  "grpc-health-port",
			Value: defaultHealthPort,
		},
		cli.IntFlag{
			Name:  "grpc-port",
			Value: defaultGrpcPort,
		},
		cli.StringFlag{
			Name:  "host",
			Value: defaultHost,
		},
	}

	app.Action = func(c *cli.Context) error {
		host := c.String("host")
		httpMetricPort := c.Int("http-metric-port")
		grpcHealthPort := c.Int("grpc-health-Port")
		grpcPort := c.Int("grpc-port")
		tls := c.Bool("tls")
		certFile := c.String("cert-file")
		keyFile := c.String("key-file")

		logrus.SetFormatter(&nested.Formatter{
			HideKeys:    true,
			FieldsOrder: []string{"group", "component", "item"},
		})

		mainLogger := logrus.WithField("group", "main")
		mainLogger.Infof("Starting %s, version %s", serviceName, version)

		if tls {
			if certFile == "" || keyFile == "" {
				err := fmt.Errorf("flag %s and %s are mandatory with flag %s", "--cert-file", "--key-file", "--tls")
				mainLogger.Errorf("Bad arguments: %s", err)
				return err
			}
			mainLogger.Infof("Using tls, keyfile: %s, certfile: %s", keyFile, certFile)
		}

		// create a cancellable context
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		// create channel for interruption
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		// defer to ensure on interrupt we have a STOP signal
		defer signal.Stop(interrupt)

		// group go routine errors
		g, ctx := errgroup.WithContext(ctx)

		// go routine for web server metrics
		g.Go(func() error {
			addr := fmt.Sprintf("%s:%d", host, httpMetricPort)
			httpServer = &http.Server{
				Addr:         addr,
				ReadTimeout:  httpReadTimeout,
				WriteTimeout: httpWriteTimeout,
			}
			mainLogger.Infof("HTTP Metrics server serving at %s", addr)

			if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
				mainLogger.Error("HTTP Metrics server: failed to listen", err)
				return err
			}

			return nil
		})

		// gRPC Health Server
		healthServer := health.NewServer()

		// go routine for health grpc server
		g.Go(func() error {
			grpcHealthServer = grpc.NewServer()
			healthPB.RegisterHealthServer(grpcHealthServer, healthServer)
			addr := fmt.Sprintf("%s:%d", host, grpcHealthPort)

			hln, err := net.Listen("tcp", addr)
			if err != nil {
				mainLogger.Errorf("gRPC Health server: %s", err)
				return err
			}
			mainLogger.Infof("gRPC health server serving at %s", addr)

			return grpcHealthServer.Serve(hln)
		})

		// gRPC server
		g.Go(func() error {
			opts := []grpc.ServerOption{
				grpc.KeepaliveEnforcementPolicy(kaep),
				grpc.KeepaliveParams(kasp),
				grpc.StatsHandler(&openCencusGrpc.ServerHandler{}),
			}

			addr := fmt.Sprintf("%s:%d", host, grpcPort)
			server, err := service.NewCarServiceImpl(serviceName, healthServer)
			if err != nil {
				return err
			}

			// We ensure that the car is gracefully stopped
			defer func() {
				if err := server.End(); err != nil {
					mainLogger.WithError(err).Error("can not end")
				}
			}()

			if tls {
				transportCredentials, err := credentials.NewServerTLSFromFile(certFile, keyFile)
				if err != nil {
					mainLogger.Errorf("Failed to generate credentials : %s", err)
					return err
				}

				opts = append(opts, grpc.Creds(transportCredentials))
			}

			ln, err := net.Listen("tcp", addr)
			if err != nil {
				mainLogger.Errorf("gRPC server: %s", err)
				return err
			}

			grpcServer = grpc.NewServer(opts...)
			grpc2.RegisterCarServiceServer(grpcServer, server)

			mainLogger.Infof("gRPC server serving at %s", addr)
			healthServer.SetServingStatus(fmt.Sprintf("grpc.health.v1.%s", serviceName), healthPB.HealthCheckResponse_SERVING)

			return grpcServer.Serve(ln)
		})

		select {
		case <-interrupt:
			break
		case <-ctx.Done():
			break
		}

		mainLogger.Warn("Received shutdown signal")

		cancel()

		healthServer.SetServingStatus(fmt.Sprintf("grpc.health.v1.%s", serviceName), healthPB.HealthCheckResponse_NOT_SERVING)

		if httpServer != nil {
			mainLogger.Warn("Shutting down HTTP Metrics server")
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer shutdownCancel()

			_ = httpServer.Shutdown(shutdownCtx)
		}

		if grpcServer != nil {
			mainLogger.Warn("Shutting down GRPC Health server")
			grpcServer.GracefulStop()
		}

		if grpcHealthServer != nil {
			mainLogger.Warn("Shutting down GRPC server")
			grpcHealthServer.GracefulStop()
		}

		err := g.Wait()
		if err != nil {
			mainLogger.Errorf("server returning an error: %s", err)
			return err
		}

		mainLogger.Infof("Application %s is shutdown", serviceName)
		mainLogger.Info("Log is disconnecting")

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
