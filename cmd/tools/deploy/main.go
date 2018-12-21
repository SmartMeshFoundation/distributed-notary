package main

import (
	"context"
	"log"

	"fmt"

	"os"

	"crypto/ecdsa"

	"github.com/SmartMeshFoundation/distributed-notary/accounts"
	ethcontracts "github.com/SmartMeshFoundation/distributed-notary/ethereum/contracts"
	"github.com/SmartMeshFoundation/distributed-notary/params"
	smccontracts "github.com/SmartMeshFoundation/distributed-notary/spectrum/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethutils "github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "address",
			Usage: "The ethereum address you would like to use and for which a keystore file exists in your local system.",
		},
		ethutils.DirectoryFlag{
			Name:  "keystore-path",
			Usage: "If you have a non-standard path for the ethereum keystore directory provide it using this argument. ",
			Value: ethutils.DirectoryString{Value: params.DefaultKeyStoreDir()},
		},
		cli.StringFlag{
			Name:  "smc-rpc-endpoint",
			Usage: `"host:port" address of Spectrum JSON-RPC server'`,
			Value: "http://127.0.0.1:8001",
		},
		cli.StringFlag{
			Name:  "eth-rpc-endpoint",
			Usage: `"host:port" address of Ethereum JSON-RPC server'`,
			Value: "http://127.0.0.1:9001",
		},
	}
	app.Action = mainctx
	app.Name = "deploy"
	app.Version = "0.1"
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func mainctx(ctx *cli.Context) error {
	address := common.HexToAddress(ctx.String("address"))
	_, keybin, err := accounts.PromptAccount(address, ctx.String("keystore-path"), "123")
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to unlock account %s", err))
	}
	fmt.Println("start to deploy ...")
	key, err := crypto.ToECDSA(keybin)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to parse priv key %s", err))
	}
	deployContractOnETH(key, createETHConn(ctx))
	deployContractOnSMC(key, createSMCConn(ctx))
	return nil
}
func deployContractOnETH(key *ecdsa.PrivateKey, conn *ethclient.Client) {
	auth := bind.NewKeyedTransactor(key)
	// Deploy Ethererum Contract
	lockedEthereumAddress, tx, _, err := ethcontracts.DeployLockedEthereum(auth, conn)
	if err != nil {
		log.Fatalf("Failed to deploy LockedEthereum contract on ethereum : %v", err)
	}
	ctx := context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	fmt.Printf("LockedEthereumAddress on ethereum = %s\n", lockedEthereumAddress.String())
	fmt.Printf("Deploy LockedEthereum contract on ethereum complete...\n")
}

func deployContractOnSMC(key *ecdsa.PrivateKey, conn *ethclient.Client) {
	auth := bind.NewKeyedTransactor(key)
	// Deploy Sepctrum Contract
	ethereumTokenAddress, tx, _, err := smccontracts.DeployEthereumToken(auth, conn)
	if err != nil {
		log.Fatalf("Failed to deploy EthereumToken contract on spectrum : %v", err)
	}
	ctx := context.Background()
	_, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	fmt.Printf("EthereumTokenAddress on spectrum = %s\n", ethereumTokenAddress.String())
	fmt.Printf("Deploy EthereumToken contract on spectrum complete...\n")
}

func createETHConn(ctx *cli.Context) *ethclient.Client {
	// Create an IPC based RPC connection to a remote node and an authorized transactor
	conn, err := ethclient.Dial(ctx.String("eth-rpc-endpoint"))
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to connect to the Ethereum client: %v", err))
	}
	return conn
}
func createSMCConn(ctx *cli.Context) *ethclient.Client {
	// Create an IPC based RPC connection to a remote node and an authorized transactor
	conn, err := ethclient.Dial(ctx.String("smc-rpc-endpoint"))
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to connect to the Spectrum client: %v", err))
	}
	return conn
}
