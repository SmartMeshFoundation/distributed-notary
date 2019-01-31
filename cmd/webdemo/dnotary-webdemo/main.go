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
			Name: "eth-rpc-point",
			Usage: `"host:port" address of ethereum   JSON-RPC server.\n'
	           'Also accepts a protocol prefix (ws:// or ipc channel) with optional port',`,
			Value: "http://127.0.0.1:8545",
		},
		cli.StringFlag{
			Name: "smc-rpc-point",
			Usage: `"host:port" address of spectrum JSON-RPC server.\n'
	           'Also accepts a protocol prefix (ws:// or ipc channel) with optional port',`,
			Value: "http://127.0.0.1:18545",
		},
		cli.IntFlag{
			Name:  "port",
			Usage: "http for service",
			Value: 8080,
		},
	}

	app.Action = mainCtx
	app.Name = "webdemo"
	app.Version = "0.1"

	err := app.Run(os.Args)
	log.Error(fmt.Sprintf("run err %s", err))

}

func mainCtx(ctx *cli.Context) {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
	rest.MainChainEndpoint = ctx.String("eth-rpc-point")
	rest.SideChainEndpoint = ctx.String("smc-rpc-point")
	rest.Port = ctx.Int("port")
	rest.RestMain()
}
