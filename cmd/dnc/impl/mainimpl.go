package dnc

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/SmartMeshFoundation/distributed-notary/accounts"
	hecoclient "github.com/SmartMeshFoundation/distributed-notary/chain/heco/client"
	hecoproxy "github.com/SmartMeshFoundation/distributed-notary/chain/heco/proxy"
	smcclient "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/client"
	smcproxy "github.com/SmartMeshFoundation/distributed-notary/chain/spectrum/proxy"
	"github.com/SmartMeshFoundation/distributed-notary/service"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/nkbai/log"
	"github.com/urfave/cli"
)

/*
	用户client,demo性质
*/

//Version version of this build
var Version string

// StartMain :
func StartMain() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		configCmd, // 管理
		/*
			lockin
		*/
		pliCmd,
		scpliCmd,
		licmd,
		cpliCmd,
		/*
			lockout
		*/
		ploCmd,
		mcploCmd,
		locmd,
		cploCmd,
		/*
			query
		*/
		queryCmd,
		/*
			test
		*/
		benchmarkCmd,
	}
	app.Name = "dnc"
	app.Version = Version
	err := app.Run(os.Args)
	if err != nil {
		os.Exit(-1)
	}
}

func getPrivateKey(addressHex, password string) (privateKey *ecdsa.PrivateKey, err error) {
	am := accounts.NewAccountManager(GlobalConfig.Keystore)
	privateKeyBin, err := am.GetPrivateKey(common.HexToAddress(addressHex), password)
	if err != nil {
		log.Error("load private key err : %s", err.Error())
		return
	}
	privateKey, err = crypto.ToECDSA(privateKeyBin)
	if err != nil {
		log.Error("load private key err : %s", err.Error())
		return
	}
	return
}

func getHecoLastBlockNumber(c *hecoclient.SafeEthClient) uint64 {
	h, err := c.HeaderByNumber(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	return h.Number.Uint64()
}

func getSmcLastBlockNumber(c *smcclient.SafeEthClient) uint64 {
	h, err := c.HeaderByNumber(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	return h.Number.Uint64()
}

func getBtcLastBlockNumber(c *rpcclient.Client) uint64 {
	_, h, err := c.GetBestBlock()
	if err != nil {
		panic(err)
	}
	return uint64(h)
}

func getSCTokenByMCName(mcName string) *service.ScTokenInfoToResponse {
	for _, sctoken := range GlobalConfig.SCTokenList {
		if sctoken.MCName == mcName {
			return &sctoken
		}
	}
	return nil
}

func getMCContractAddressByMCName(mcName string) common.Address {
	if GlobalConfig.SCTokenList == nil || len(GlobalConfig.SCTokenList) == 0 {
		fmt.Println("must run dnc config refresh first")
		os.Exit(-1)
	}
	for _, sctoken := range GlobalConfig.SCTokenList {
		if sctoken.MCName == mcName {
			return sctoken.MCLockedContractAddress
		}
	}
	fmt.Printf("can not found mc contract address of %s\n", mcName)
	os.Exit(-1)
	return utils.EmptyAddress
}

func getSCContractAddressByMCName(mcName string) common.Address {
	if GlobalConfig.SCTokenList == nil || len(GlobalConfig.SCTokenList) == 0 {
		fmt.Println("must run dnc config refresh first")
		os.Exit(-1)
	}
	for _, sctoken := range GlobalConfig.SCTokenList {
		if sctoken.MCName == mcName {
			return sctoken.SCToken
		}
	}
	fmt.Printf("can not found sc token address of %s\n", mcName)
	os.Exit(-1)
	return utils.EmptyAddress
}

func eth2Wei(ethAmount int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(params.Shannon)), big.NewInt(ethAmount))
}

func wei2Eth(weiAmount *big.Int) *big.Int {
	return new(big.Int).Div(weiAmount, big.NewInt(int64(params.Shannon)))
}

func getHecoConn() *hecoclient.SafeEthClient {
	ctx2, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	c, err := ethclient.DialContext(ctx2, GlobalConfig.HecoRPCEndpoint)
	cancelFunc()
	if err != nil {
		fmt.Println("connect to eth fail : ", err)
		os.Exit(-1)
	}
	return hecoclient.NewSafeClient(c)
}
func getSmcConn() *smcclient.SafeEthClient {
	ctx2, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	c, err := ethclient.DialContext(ctx2, GlobalConfig.SmcRPCEndpoint)
	cancelFunc()
	if err != nil {
		fmt.Println("connect to eth fail : ", err)
		os.Exit(-1)
	}
	return smcclient.NewSafeClient(c)
}

func getSCContractProxy(mcName string) (*hecoclient.SafeEthClient, *hecoproxy.SideChainErc20TokenProxy) {
	if GlobalConfig.SCTokenList == nil {
		fmt.Println("must run dnc config refresh first")
		os.Exit(-1)
	}
	// 1. init connect
	ctx2, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	c, err := ethclient.DialContext(ctx2, GlobalConfig.HecoRPCEndpoint)
	cancelFunc()
	if err != nil {
		fmt.Println("connect to eth fail : ", err)
		os.Exit(-1)
	}
	conn := hecoclient.NewSafeClient(c)

	lastBlockNumber, err := conn.HeaderByNumber(context.Background(), nil)
	if err != nil {
		fmt.Println("HeaderByNumber err : ", err)
		os.Exit(-1)
	}
	fmt.Printf("[SC] lasted block number = %d\n", lastBlockNumber.Number.Uint64())

	// 2. init contract proxy
	cp, err := hecoproxy.NewSideChainErc20TokenProxy(conn, getSCContractAddressByMCName(mcName))
	if err != nil {
		fmt.Println("init contract proxy err : ", err)
		os.Exit(-1)
	}
	return conn, cp
}

func getMCContractProxy(mcName string) (*smcclient.SafeEthClient, *smcproxy.LockedSpectrumProxy) {
	if GlobalConfig.SCTokenList == nil {
		fmt.Println("must run dnc config refresh first")
		os.Exit(-1)
	}
	// 1. init connect
	ctx2, cancelFunc := context.WithTimeout(context.Background(), 3*time.Second)
	c, err := ethclient.DialContext(ctx2, GlobalConfig.SmcRPCEndpoint)
	cancelFunc()
	if err != nil {
		fmt.Println("connect to eth fail : ", err)
		os.Exit(-1)
	}
	conn := smcclient.NewSafeClient(c)

	lastBlockNumber, err := conn.HeaderByNumber(context.Background(), nil)
	if err != nil {
		fmt.Println("HeaderByNumber err : ", err)
		os.Exit(-1)
	}
	fmt.Printf("[MC] lasted block number = %d\n", lastBlockNumber.Number.Uint64())
	// 2. init contract proxy
	cp, err := smcproxy.NewLockedSpectrumProxy(conn, getMCContractAddressByMCName(mcName))
	if err != nil {
		fmt.Println("init contract proxy err : ", err)
		os.Exit(-1)
	}
	return conn, cp
}
