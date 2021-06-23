package dnc

import (
	"context"
	"crypto/ecdsa"
	"net/http"

	"fmt"
	"os"

	"time"

	"sync"

	"math/rand"

	"sort"

	"github.com/SmartMeshFoundation/distributed-notary/api"
	"github.com/SmartMeshFoundation/distributed-notary/api/userapi"
	"github.com/SmartMeshFoundation/distributed-notary/chain/spectrum"
	"github.com/SmartMeshFoundation/distributed-notary/models"
	"github.com/SmartMeshFoundation/distributed-notary/service"
	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli"
)

var benchmarkCmd = cli.Command{
	Name:      "benchmark",
	ShortName: "bm",
	Usage:     "run benchmark test",
	Action:    benchmark,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "mcname",
			Usage: "name of main chain contract which you want to lockin",
			Value: "ethereum",
		},
		cli.BoolFlag{
			Name:  "pkn",
			Usage: "run benchmark of pkn ",
		},
		cli.BoolFlag{
			Name:  "dsm",
			Usage: "run benchmark of dsm ",
		},
		cli.IntFlag{
			Name:  "num",
			Usage: "num of gorouting",
			Value: 100,
		},
	},
}

func benchmark(ctx *cli.Context) {
	mcname := ctx.String("mcname")
	num := ctx.Int("num")
	if ctx.Bool("pkn") {
		benchmarkPKN(mcname, num)
		return
	}
	if ctx.Bool("dsm") {
		benchmarkDSM(mcname, num)
		return
	}
	fmt.Println("need --pkn or --dsm")
	os.Exit(-1)
}

func benchmarkPKN(mcName string, num int) {
	fmt.Printf("==> PKN Benchmark  START ...\n")
	var used []float64
	var usedLock sync.Mutex
	start := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			url := fmt.Sprintf("http://127.0.0.1:803%d/api/1/admin/private-key", rand.Intn(7))
			var resp api.BaseResponse
			s := time.Now()
			err2 := call(http.MethodPut, url, "", &resp)
			use := time.Since(s)
			if err2 != nil {
				fmt.Printf("call %s err :%s", url, err2.Error())
				os.Exit(-1)
			}
			var pk service.PrivateKeyInfoToResponse
			err2 = resp.ParseData(&pk)
			if err2 != nil {
				fmt.Printf("parse response data =%s err :%s", resp.Data, err2.Error())
				os.Exit(-1)
			}
			if pk.Status != models.PrivateKeyNegotiateStatusFinished {
				fmt.Printf("private key status fail :\n%s", utils.ToJSONStringFormat(pk))
			}
			usedLock.Lock()
			used = append(used, use.Seconds())
			usedLock.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	total := time.Since(start)
	// 5. 输出
	sort.Float64s(used)
	fmt.Printf("==> PKN Benchmark  END SUCCESS\n")
	fmt.Printf("==> total use %f seconds, avg %f seconds\n", total.Seconds(), total.Seconds()/float64(num))
	if num <= 20 {
		for index, use := range used {
			fmt.Printf("%02d : %f\n", index+1, use)
		}
		return
	}
	// 前10
	fmt.Println("==> fastest 10 :")
	for i := 0; i < 10; i++ {
		fmt.Printf("%04d : %f\n", i+1, used[i])
	}
	// 后10
	fmt.Println("==> slowest 10 :")
	for i := num - 10; i < num; i++ {
		fmt.Printf("%04d : %f\n", i+1, used[i])
	}
}

func benchmarkDSM(mcName string, num int) {
	if GlobalConfig.SCTokenList == nil {
		fmt.Println("must run dnc config refresh first")
		os.Exit(-1)
	}
	fmt.Printf("==> DSM Benchmark prepare start mcName=%s num=%d\n", mcName, num)
	mcKey, err := getPrivateKey(GlobalConfig.SmcUserAddress, GlobalConfig.SmcUserPassword)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	scToken := getSCTokenByMCName(mcName)
	_, mcp := getMCContractProxy(mcName)
	es, err := spectrum.NewSMCService(GlobalConfig.SmcRPCEndpoint)
	if err != nil {
		fmt.Println("connect to eth fail : ", err)
		os.Exit(-1)
	}
	//1. 生成100个私钥,100个密码
	var keys []*ecdsa.PrivateKey
	var secretHashs []common.Hash
	for i := 0; i < num; i++ {
		k, err2 := crypto.GenerateKey()
		if err2 != nil {
			panic(err2)
		}
		keys = append(keys, k)
		secret := utils.NewRandomHash()
		secretHashs = append(secretHashs, utils.ShaSecret(secret[:]))
	}

	amount := int64(1)
	expiration := uint64(10000)
	expiration2 := getEthLastBlockNumber(es.GetClient()) + expiration
	fmt.Printf("==> wait for init accounts , include transfer money and PLI  ...\n")
	baseNonce, err := es.GetClient().NonceAt(context.Background(), crypto.PubkeyToAddress(mcKey.PublicKey), nil)
	if err != nil {
		fmt.Println("GetClient err ", err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(len(keys))
	for i, ckey := range keys {
		go func(index int, key *ecdsa.PrivateKey, nonce int) {
			account := crypto.PubkeyToAddress(key.PublicKey)
			secretHash := secretHashs[index]
			//2. 主链转账1eth
			err := es.Transfer10ToAccount(mcKey, account, eth2Wei(10*1000000000), nonce)
			if err != nil {
				fmt.Println("transfer eth to account fail : ", err)
				os.Exit(-1)
			}
			//3. 调用合约pli
			auth := bind.NewKeyedTransactor(key)
			err = mcp.PrepareLockin(auth, account.String(), secretHash, expiration2, eth2Wei(amount))
			if err != nil {
				fmt.Println("pli fail : ", err)
				os.Exit(-1)
			}
			wg.Done()
		}(i, ckey, int(baseNonce)+i)

	}
	wg.Wait()
	// 3.5 等待pli事件完成
	fmt.Printf("==> wait for notary to confirm pli event ...\n")
	time.Sleep(20 * time.Second)
	// 4. 构造scpli请求,  每个私钥随机选取一个公证人调用scpli
	reqs := getBenchmarkRequests(scToken.SCToken, keys, secretHashs)
	// 5. 调用
	fmt.Printf("==> DSM Benchmark  START ...\n")
	start := time.Now()
	wg = sync.WaitGroup{}
	wg.Add(num)
	cnt := 0
	for _, r := range reqs {
		go func(r *req) {
			var resp api.BaseResponse
			err2 := call(http.MethodPost, r.url, r.payload, &resp)
			if err2 != nil {
				fmt.Printf("call %s with payload=%s err :%s", r.url, r.payload, err2.Error())
				os.Exit(-1)
			}
			wg.Done()
			cnt++
			fmt.Printf("dnc cnt=%d\n", cnt)
		}(r)
	}
	wg.Wait()
	total := time.Since(start)
	// 5. 输出
	fmt.Printf("==> DSM Benchmark  END SUCCESS\n")
	fmt.Printf("==> total use %f seconds, avg %f seconds\n", total.Seconds(), total.Seconds()/float64(num))
}

type req struct {
	url     string
	payload string
}

func getBenchmarkRequests(scTokenAddress common.Address, keys []*ecdsa.PrivateKey, secretHashs []common.Hash) (requests []*req) {
	for index, key := range keys {
		address := crypto.PubkeyToAddress(key.PublicKey)
		body := &userapi.SCPrepareLockinRequest{
			BaseReq:              api.NewBaseReq(userapi.APIUserNameSCPrepareLockin),
			BaseReqWithResponse:  api.NewBaseReqWithResponse(),
			BaseReqWithSCToken:   api.NewBaseReqWithSCToken(scTokenAddress),
			BaseReqWithSignature: api.NewBaseReqWithSignature(),
			SecretHash:           secretHashs[index],
			MCUserAddress:        address[:],
			//SCUserAddress:        address,
		}
		body.Sign(body, key)
		requests = append(requests, &req{
			url:     fmt.Sprintf("http://127.0.0.1:803%d%s", 0, "/api/1/user/scpreparelockin/"+scTokenAddress.String()),
			payload: utils.ToJSONString(body),
		})
	}
	return
}
