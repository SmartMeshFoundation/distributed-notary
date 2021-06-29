package main

import (
	"os"

	"github.com/SmartMeshFoundation/distributed-notary/cmd/nonce_server/nonceapi"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/urfave/cli"
)

//Version version of this build
var Version = "v0.1"

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "datadir",
			Usage: "Directory for storing nonce-server data.",
			Value: "./.nonce-server-data",
		},
		cli.StringFlag{
			Name:  "listen",
			Usage: "host:port of nsAPI listen",
			Value: "0.0.0.0:8020",
		},
		cli.StringFlag{
			Name:  "smc-rpc-endpoint",
			Usage: "host:port of spectrum rpc server",
			Value: "http://106.52.171.12:18003",
		},
		cli.StringFlag{
			Name:  "heco-rpc-endpoint",
			Usage: "host:port of heco rpc server",
			Value: "http://106.52.171.12:12001",
		},
	}
	app.Action = StartMain
	app.Name = "nonce_server"
	app.Version = Version
	err := app.Run(os.Args)
	if err != nil {
		os.Exit(-1)
	}
}

// StartMain :
func StartMain(ctx *cli.Context) {
	dataDir := ctx.String("datadir")
	host := ctx.String("listen")
	smcRPCEndPoint := ctx.String("smc-rpc-endpoint")
	hecoRPCEndPoint := ctx.String("heco-rpc-endpoint")
	// 1. 打开db
	db := models.SetUpDB("sqlite3", dataDir)
	// 2. init nsAPI
	api := nonceapi.NewNonceServerAPI(host)
	// 3. init service
	service := newNonceService(db, api, smcRPCEndPoint, hecoRPCEndPoint)
	go service.start()
	// 2. 启动api
	api.Start(true)
}
