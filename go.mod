module github.com/iam-merlin/carlos

go 1.13

require (
	github.com/akhenakh/ocgrpc_propagation v0.0.0-20190306172630-8bd08bcc1ad4
	github.com/antonfisher/nested-logrus-formatter v1.0.2
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.22.1
	go.opencensus.io v0.22.2 // indirect
	gobot.io/x/gobot v0.0.0-00010101000000-000000000000
	golang.org/x/net v0.0.0-20191109021931-daa7c04131f5 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys v0.0.0-20191105231009-c1f44814a5cd // indirect
	google.golang.org/genproto v0.0.0-20191108220845-16a3f7862a1a // indirect
	google.golang.org/grpc v1.25.1
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.5 // indirect
)

replace gobot.io/x/gobot => github.com/iam-merlin/gobot v1.14.1-0.20191106204422-2dfe9651f25f

replace gobot.io/x/gobot/drivers/i2c => github.com/iam-merlin/gobot/drivers/i2c v1.14.1-0.20191106204422-2dfe9651f25f

replace gobot.io/x/gobot/platforms/raspi => github.com/iam-merlin/gobot/platforms/raspi v1.14.1-0.20191106204422-2dfe9651f25f
