package main

import (
	"fmt"
	"os"

	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/SmartMeshFoundation/distributed-notary/cmd/webdemo/rest"

	"github.com/SmartMeshFoundation/Photon/log"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "heco-rpc-point",
			Usage: `"host:port" address of heco JSON-RPC server.\n'
	           'Also accepts a protocol prefix (ws:// or ipc channel) with optional port',`,
			Value: "http://106.52.171.12:12001",
		},
		cli.StringFlag{
			Name: "smc-rpc-point",
			Usage: `"host:port" address of spectrum JSON-RPC server.\n'
	           'Also accepts a protocol prefix (ws:// or ipc channel) with optional port',`,
			Value: "http://106.52.171.12:18003",
		},
		cli.IntFlag{
			Name:  "port",
			Usage: "http for service",
			Value: 8081,
		},
	}

	app.Action = mainCtx
	app.Name = "webdemo"
	app.Version = "0.2"

	err := app.Run(os.Args)
	log.Error(fmt.Sprintf("run err %s", err))

}

func mainCtx(ctx *cli.Context) {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
	rest.MainChainEndpoint = ctx.String("smc-rpc-point")
	rest.SideChainEndpoint = ctx.String("heco-rpc-point")
	rest.Port = ctx.Int("port")
	rest.RestMain()
}
