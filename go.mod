module github.com/iam-merlin/carlos

go 1.13

require (
	cloud.google.com/go v0.47.0 // indirect
	cloud.google.com/go/storage v1.2.1 // indirect
	code.gitea.io/sdk/gitea v0.0.0-20191106151626-e4082d89cc3b // indirect
	contrib.go.opencensus.io/exporter/ocagent v0.6.0 // indirect
	github.com/Azure/azure-pipeline-go v0.2.2 // indirect
	github.com/Azure/azure-sdk-for-go v36.1.0+incompatible // indirect
	github.com/Azure/azure-storage-blob-go v0.8.0 // indirect
	github.com/Azure/go-autorest v13.3.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/adal v0.8.0 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/akhenakh/ocgrpc_propagation v0.0.0-20190306172630-8bd08bcc1ad4
	github.com/antonfisher/nested-logrus-formatter v1.0.2
	github.com/aws/aws-sdk-go v1.25.31 // indirect
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/goreleaser/goreleaser v0.120.6 // indirect
	github.com/goreleaser/nfpm v1.1.5 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.12.1 // indirect
	github.com/jstemmer/go-junit-report v0.9.1 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/mattn/go-ieproxy v0.0.0-20190805055040-f9202b1cfdeb // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.22.1
	github.com/xanzy/go-gitlab v0.22.0 // indirect
	go.opencensus.io v0.22.2 // indirect
	gobot.io/x/gobot v1.14.1
	golang.org/x/crypto v0.0.0-20191108234033-bd318be0434a // indirect
	golang.org/x/net v0.0.0-20191109021931-daa7c04131f5 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys v0.0.0-20191105231009-c1f44814a5cd // indirect
	golang.org/x/tools v0.0.0-20191108193012-7d206e10da11 // indirect
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 // indirect
	google.golang.org/genproto v0.0.0-20191108220845-16a3f7862a1a // indirect
	google.golang.org/grpc v1.25.1
)

replace gobot.io/x/gobot => github.com/iam-merlin/gobot v1.14.1-0.20191106204422-2dfe9651f25f

replace gobot.io/x/gobot/drivers/i2c => github.com/iam-merlin/gobot/drivers/i2c v1.14.1-0.20191106204422-2dfe9651f25f

replace gobot.io/x/gobot/platforms/raspi => github.com/iam-merlin/gobot/platforms/raspi v1.14.1-0.20191106204422-2dfe9651f25f
