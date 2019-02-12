package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"encoding/hex"
	"path/filepath"

	"github.com/SmartMeshFoundation/distributed-notary/params"
	"github.com/SmartMeshFoundation/distributed-notary/service"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

//Version version of this build
var Version string

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "datadir",
			Usage: "Directory for storing distributed-notary data.",
			Value: "./.dnotary-data",
		},
		cli.StringFlag{
			Name:  "address",
			Usage: "The ethereum address you would like distributed-notary to use and for which a keystore file exists in your local system.",
		},
		cli.StringFlag{
			Name:  "keystore-path",
			Usage: "If you have a non-standard path for the ethereum keystore directory provide it using this argument. ",
			Value: "./keystore",
		},
		cli.StringFlag{
			Name:  "password-file",
			Usage: "Text file containing password for provided account",
			Value: "123",
		},
		cli.StringFlag{
			Name:  "notary-config-file",
			Usage: "config file containing notary info list. only need when first start",
			Value: "./notary.conf",
		},
		cli.StringFlag{
			Name:  "smc-rpc-point",
			Usage: "host:port of spectrum rpc server",
			Value: "http://127.0.0.1:17888",
		},
		cli.StringFlag{
			Name:  "eth-rpc-point",
			Usage: "host:port of spectrum rpc server",
			Value: "http://127.0.0.1:19888",
		},
		cli.StringFlag{
			Name:  "user-listen",
			Usage: "host:port of user api listen",
			Value: "127.0.0.1:3330",
		},
		cli.StringFlag{
			Name:  "notary-listen",
			Usage: "host:port of notary api listen",
			Value: "127.0.0.1:33300",
		},
		cli.StringFlag{
			Name:  "nonce-server",
			Usage: "http://host:port of nonce server",
			Value: "http://127.0.0.1:8020",
		},
	}
	app.Action = startMain
	app.Name = "distributed-notary"
	app.Version = Version
	err := app.Run(os.Args)
	if err != nil {
		os.Exit(-1)
	}
}

func startMain(ctx *cli.Context) {
	err := mainCtx(ctx)
	if err != nil {
		fmt.Printf("quit with err : %s\n", err.Error())
		os.Exit(-1)
	}
}

func mainCtx(ctx *cli.Context) (err error) {
	// 1. 加载配置
	var cfg *params.Config
	cfg, err = config(ctx)
	if err != nil {
		return
	}
	// 2. 初始化DispatchService
	ds, err := service.NewDispatchService(cfg)
	if err != nil {
		return
	}
	// 3. DispatchService.Start()
	err = ds.Start()
	if err != nil {
		return
	}
	return
}

func config(ctx *cli.Context) (cfg *params.Config, err error) {
	cfg = &params.Config{}
	// 1. address
	cfg.Address = common.HexToAddress(ctx.String("address"))
	if cfg.Address == utils.EmptyAddress {
		err = fmt.Errorf("can not start without --address")
		return
	}
	// 2. datadir
	datadir := ctx.String("datadir")
	if !utils.Exists(datadir) {
		err = os.MkdirAll(datadir, os.ModePerm)
		if err != nil {
			err = fmt.Errorf("datadir:%s doesn't exist and cannot create %v", datadir, err)
			return
		}
	}
	userDbPath := hex.EncodeToString(cfg.Address[:])
	userDbPath = userDbPath[:8]
	userDbPath = filepath.Join(datadir, userDbPath)
	if !utils.Exists(userDbPath) {
		err = os.MkdirAll(userDbPath, os.ModePerm)
		if err != nil {
			err = fmt.Errorf("datadir:%s doesn't exist and cannot create %v", userDbPath, err)
			return
		}
	}
	databasePath := filepath.Join(userDbPath, "log.db")
	cfg.DataBasePath = databasePath
	// 3. keystore-path
	cfg.KeystorePath = ctx.String("keystore-path")
	// 4. password-file
	passwordFile := ctx.String("password-file")
	//#nosec
	data, err2 := ioutil.ReadFile(passwordFile)
	if err2 != nil {
		data = []byte(passwordFile)
	}
	cfg.Password = string(data)
	// 5. notary-config-file
	cfg.NotaryConfFilePath = ctx.String("notary-config-file")
	// 6. smc-rpc-point
	cfg.SmcRPCEndPoint = ctx.String("smc-rpc-point")
	// 7. smc-rpc-point
	cfg.EthRPCEndPoint = ctx.String("eth-rpc-point")
	// 8. user-listen
	cfg.UserAPIListen = ctx.String("user-listen")
	// 9. notary-listen
	cfg.NotaryAPIListen = ctx.String("notary-listen")
	// 10. nonce-server
	cfg.NonceServerHost = ctx.String("nonce-server")
	return
}
