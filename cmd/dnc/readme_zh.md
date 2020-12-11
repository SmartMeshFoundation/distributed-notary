
# dnc 使用说明

## 1. 初始化
### 1.1 dnc 使用帮助
```
 ./dnc --help
NAME:
   dnc - A new cli application

USAGE:
   dnc [global options] command [command options] [arguments...]

COMMANDS:
     config, c                           manage all config of dnc
     prepare-lock-in, pli                call main chain contract prepare lock in
     side-chain-prepare-lock-in, scpli   call SCPrepareLockin API of notary
     lock-in, li                         call side chain contract lock in
     cancel-prepare-lock-in, cpli        call main chain contract cancel prepare lock in
     prepare-lock-out, plo               call side chain contract prepare lock out
     main-chain-prepare-lock-out, mcplo  call MCPrepareLockout API of notary
     lock-out, lo                        call main chain contract lock out
     cancel-prepare-lock-out, cplo       call side chain contract cancel prepare lock out
     query, q                            query lockin/lockout info on sc/mc
     benchmark, bm                       run benchmark test
     help, h                             Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
### 1.2 dnc 配置
请使用附带的dnc.json文件,请忽略其中有关btc的任何内容

### 1.3 获取侧链token合约地址
```
./dnc config refresh
```
这时候打开dnc.json可以看到:
```
{
	"notary_host": "http://transport01.smartmesh.cn:8032",
	"keystore": "./keystore",
	"btc_rpc_user": "wuhan",
	"btc_rpc_pass": "wuhan",
	"btc_rpc_cert_file_path": "",
	"btc_rpc_endpoint": "192.168.124.13:18556",
	"btc_user_address": "SgEQfVdPqBS65jpSNLoddAa9kCouqqxGrY",
	"btc_wallet_rpc_cert_file_path": "",
	"btc_wallet_rpc_endpoint": "192.168.124.13:18554",
	"eth_user_address": "0x201b20123b3c489b47fde27ce5b451a0fa55fd60",
	"eth_user_password": "123",
	"eth_rpc_endpoint": "http://106.52.171.12:18003",
	"smc_user_address": "0x201b20123b3c489b47fde27ce5b451a0fa55fd60",
	"smc_user_password": "123",
	"smc_rpc_endpoint": "http://106.52.171.12:17004",
	"sc_token_list": [
		{
			"sc_token": "0x7b6cfc6a9edd6a1f81af4f8b860eebc26bf1ae4b",
			"sc_token_name": "ethereum-Token",
			"sc_token_owner_key": "0xf16ebdf72d336731f7f09dcde52637ee7d16d9af2c8ed68777f8582957dd73a1",
			"mc_locked_contract_address": "0x7b6cfc6a9edd6a1f81af4f8b860eebc26bf1ae4b",
			"mc_name": "ethereum",
			"mc_locked_contract_owner_key": "0xf16ebdf72d336731f7f09dcde52637ee7d16d9af2c8ed68777f8582957dd73a1",
			"create_time": "2020-12-11 17:32:51 +0800 CST",
			"organiser_id": 2
		}
	],
	"run_time": null
}
```
sc_token_list 中的sc_token`0x7b6cfc6a9edd6a1f81af4f8b860eebc26bf1ae4b` 就是我们要使用的合约地址.

#### 1.4 切换账号
请在keystore中放置相应的文件,并且修改1.3节中的`smc_user_address`和`eth_user_address` 以及相应的密码.
请勿使用`dnc config reset`命令来重置该配置文件.

#### 1.5 查询合约当前状态

```
./dnc q --all
===> MC/SC Lasted BlockNumber info :
[MC] lasted block number = 9238490
[SC] lasted block number = 3674516

===> MC/SC User account info :
[MC]user 201b2012 account balance : 997138423763722556
[SC]user 201b2012 sctoken balance : 8000

===> MC/SC Contract account info :
[MC]contract 7b6cfc6a account balance : 8040

===> MC/SC Contract data info :
[MC]data of lockin  :
	 account    =  0x201b20123b3c489b47fde27ce5b451a0fa55fd60
	 secretHash =  0xcb8887f75fe79e1fe0f2a2588cad425f9fef2d2eb1315b40b718868678cc10c8
	 expiration =  9239383
	 amount     =  20
[SC]data of lockin  : Empty
[MC]data of lockout : Empty
[SC]data of lockout : Empty
```
## 2. 主链到侧链跨链

### 2.1 主链PrepareLockIn

```
./dnc pli --amount 20
```
这里的amount是20wei,以wei为单位

成功后得到如下输出:
```
start to prepare lockin :
 ======> [chain=ethereum amount=20 expiartion=900]
 ======> [secret=0xf1b820fd53c8747bd959f7172d5d38fe5f0174d1cfc15a7db3205b83584dadd2, secretHash=0x474ce39fa753dbe30fe4d43e37a0959b2713a291a22053a88981b37ad7d90be1]
PrepareLockin tx=
	TX(e811c7719f51c1db85d0f53c3fbde0480db8940f8520e676a1b386eef44781d3)
	Contract: false
	From:     201b20123b3c489b47fde27ce5b451a0fa55fd60
	To:       7b6cfc6a9edd6a1f81af4f8b860eebc26bf1ae4b
	Nonce:    7
	GasPrice: 0x11520a2b8
	GasLimit  0x14a60
	Value:    0x14
	Data:     0xe0ae1a81474ce39fa753dbe30fe4d43e37a0959b2713a291a22053a88981b37ad7d90be100000000000000000000000000000000000000000000000000000000008cf5f4
	V:        0x2a
	R:        0x43217c3913b5bf30e7434ea911fbf935b829bbda9c3221bf44bcc8ec934e52a8
	S:        0xa11de168ba80a6a1fff789503101624304f8cfbf0e2ee48cde1e4f4f39e33f9
	Hex:      f8aa0785011520a2b883014a60947b6cfc6a9edd6a1f81af4f8b860eebc26bf1ae4b14b844e0ae1a81474ce39fa753dbe30fe4d43e37a0959b2713a291a22053a88981b37ad7d90be100000000000000000000000000000000000000000000000000000000008cf5f42aa043217c3913b5bf30e7434ea911fbf935b829bbda9c3221bf44bcc8ec934e52a8a00a11de168ba80a6a1fff789503101624304f8cfbf0e2ee48cde1e4f4f39e33f9

PrepareLockin on ethereum SUCCESS
```
这次跨链使用的
- 密码:0xf1b820fd53c8747bd959f7172d5d38fe5f0174d1cfc15a7db3205b83584dadd2 这个是私有的,目前不会在链上公开
- 密码hash: 0x474ce39fa753dbe30fe4d43e37a0959b2713a291a22053a88981b37ad7d90be1
- 跨链的金额: 20
- 过期时间: 900块以后 
- Tx: [e811c7719f51c1db85d0f53c3fbde0480db8940f8520e676a1b386eef44781d3](https://ropsten.etherscan.io/tx/0xe811c7719f51c1db85d0f53c3fbde0480db8940f8520e676a1b386eef44781d3)

### 2.2 通知公证人

通知公证人我刚刚进行了跨链行为,需要公证人在侧链进行preparelockin
```
 ./dnc scpli
```

得到如下输出:
```
req={"name":"User-SCPrepareLockin","request_id":"0xac97529f278d985a9be69cd6013cc52d038212a1ce703c24a536dd2480dc6c45","sc_token_address":"0x7b6cfc6a9edd6a1f81af4f8b860eebc26bf1ae4b","signer":"03f2307da626a34b0ebd1732b97df2517630bf699198c43b54f1d1b7696b9edd22","secret_hash":"0x474ce39fa753dbe30fe4d43e37a0959b2713a291a22053a88981b37ad7d90be1","mc_user_address":"0x201b20123b3c489b47fde27ce5b451a0fa55fd60"}
digest=427d86fc0c453fcae33a00ec7a23f9614a3fca94d1592c896a7111cd7e0d0c85
SCPrepareLockin SUCCESS
```
req是向公证人提交的数据 详细见说明文档.

### 2.3 侧链lockin获得erc20 代币
通过在侧链上调用lockin来公布密码`0xf1b820fd53c8747bd959f7172d5d38fe5f0174d1cfc15a7db3205b83584dadd2`
```
./dnc li
```

```
[SC] lasted block number = 3673719
start to lockin :
 ======> [account=0x201b20123b3c489b47fde27ce5b451a0fa55fd60 secret=0xf1b820fd53c8747bd959f7172d5d38fe5f0174d1cfc15a7db3205b83584dadd2 secretHash=0x474ce39fa753dbe30fe4d43e37a0959b2713a291a22053a88981b37ad7d90be1]
lockin tx=
	TX(7dc1f7b614f9432bc892d5243657a9853101b2e6735a22081f5695643a7abb0c)
	Contract: false
	From:     201b20123b3c489b47fde27ce5b451a0fa55fd60
	To:       7b6cfc6a9edd6a1f81af4f8b860eebc26bf1ae4b
	Nonce:    13
	GasPrice: 0x430e23400
	GasLimit  0x110f2
	Value:    0x0
	Data:     0x7fd408d2000000000000000000000000201b20123b3c489b47fde27ce5b451a0fa55fd60f1b820fd53c8747bd959f7172d5d38fe5f0174d1cfc15a7db3205b83584dadd2
	V:        0x2a
	R:        0x983032017b541a0d3f06de1fcaff02b6d23264d0d190aa6bf221f59f9c6c8e3c
	S:        0x755f93079170cdc8362418e561a8a8197a5116753cef43a8826537749b5e612a
	Hex:      f8aa0d850430e23400830110f2947b6cfc6a9edd6a1f81af4f8b860eebc26bf1ae4b80b8447fd408d2000000000000000000000000201b20123b3c489b47fde27ce5b451a0fa55fd60f1b820fd53c8747bd959f7172d5d38fe5f0174d1cfc15a7db3205b83584dadd22aa0983032017b541a0d3f06de1fcaff02b6d23264d0d190aa6bf221f59f9c6c8e3ca0755f93079170cdc8362418e561a8a8197a5116753cef43a8826537749b5e612a

Lockin SUCCESS
```
其中:
- 公布的密码: 0xf1b820fd53c8747bd959f7172d5d38fe5f0174d1cfc15a7db3205b83584dadd2
- Tx [7dc1f7b614f9432bc892d5243657a9853101b2e6735a22081f5695643a7abb0c](https://chain.smartmesh.cn/tx.html?hash=0x7dc1f7b614f9432bc892d5243657a9853101b2e6735a22081f5695643a7abb0c)

#### 2.4 查询结果
到此,从主链到侧链完成,可以看到
```
./dnc q --all 
```

```
===> MC/SC Lasted BlockNumber info :
[MC] lasted block number = 9237129
[SC] lasted block number = 3673728

===> MC/SC User account info :
[MC]user 201b2012 account balance : 998550845489920316
[SC]user 201b2012 sctoken balance : 20

===> MC/SC Contract account info :
[MC]contract 7b6cfc6a account balance : 20

===> MC/SC Contract data info :
[MC]data of lockin  : Empty
[SC]data of lockin  : Empty
[MC]data of lockout : Empty
[SC]data of lockout : Empty
```

`[SC]user 201b2012 sctoken balance : 20` 说明account `201b2012` 在侧链上有了20wei的代币
`[MC]contract 7b6cfc6a account balance : 20` 说明主链合约上锁定了20wei的eth

## 3. 侧链到主链跨链
侧链跨回主链,就是将刚刚得到的eth代币换回eth.

### 3.1 侧链preparelockout 
```
./dnc plo --amount 20
```

得到输出:
```
start to prepare lockout :
 ======> [contract=0x7b6Cfc6A9eDD6A1f81aF4f8b860eEbc26Bf1ae4b amount=20 expiartion=1000]
 ======> [secret=0x8bb14feefdcdf07fdbda4297be7e64f4494db5af895f9f0eceb39616bdf5c2b3, secretHash=0x40ddb2475d5b7a74ba9f28f6701e355c570f4e3d3fb965f352c67c17edd5e2b3]
[SC] lasted block number = 3673731
 ======> [token balance=20]
call params :
callerAddress =  0x201B20123b3C489b47Fde27ce5b451a0fA55FD60
secretHash    =  0x40ddb2475d5b7a74ba9f28f6701e355c570f4e3d3fb965f352c67c17edd5e2b3
expiration    =  3674731
amount        =  20
preparelockout tx=
	TX(c758dfab62b4654abf11b45395ae84f93ca91487be1d202b4654a716cc123843)
	Contract: false
	From:     201b20123b3c489b47fde27ce5b451a0fa55fd60
	To:       7b6cfc6a9edd6a1f81af4f8b860eebc26bf1ae4b
	Nonce:    14
	GasPrice: 0x430e23400
	GasLimit  0x16650
	Value:    0x0
	Data:     0x92d062cd40ddb2475d5b7a74ba9f28f6701e355c570f4e3d3fb965f352c67c17edd5e2b3000000000000000000000000000000000000000000000000000000000038126b0000000000000000000000000000000000000000000000000000000000000014
	V:        0x29
	R:        0x25aa21ad26eeb50e2d872f16079b046ee4967ea71936d7d5786cb1fae1a02e6c
	S:        0x7350cbbfdc177c3c1b81bf4846ad73ea243b20de857d1c396029738bda2d2b7b
	Hex:      f8ca0e850430e2340083016650947b6cfc6a9edd6a1f81af4f8b860eebc26bf1ae4b80b86492d062cd40ddb2475d5b7a74ba9f28f6701e355c570f4e3d3fb965f352c67c17edd5e2b3000000000000000000000000000000000000000000000000000000000038126b000000000000000000000000000000000000000000000000000000000000001429a025aa21ad26eeb50e2d872f16079b046ee4967ea71936d7d5786cb1fae1a02e6ca07350cbbfdc177c3c1b81bf4846ad73ea243b20de857d1c396029738bda2d2b7b

PrepareLockout SUCCESS
```
这次跨链使用的:
- 密码: 0x78f6c1c7dd75b8e34bdce6496a368be86ab556c3579e49a9e532cf534abbef5c
- 密码hash: 0xcede52a047a88b6c133ac79f83a606ce8e854a40af52eac4d153bbe315d0feb4
- 过期相对块数: 1000
- 金额: 20wei
- Tx [c758dfab62b4654abf11b45395ae84f93ca91487be1d202b4654a716cc123843](https://chain.smartmesh.cn/tx.html?hash=0xc758dfab62b4654abf11b45395ae84f93ca91487be1d202b4654a716cc123843)

### 3.2 通知公证人

```
./dnc mcplo
```
得到输出 :
```
req={"name":"User-MCPrepareLockout","request_id":"0x05a4cdd2d917ec48dce718ddc4b05425eeb8290ca9a9023e0c786a79dd8fdad1","sc_token_address":"0x7b6cfc6a9edd6a1f81af4f8b860eebc26bf1ae4b","signer":"03f2307da626a34b0ebd1732b97df2517630bf699198c43b54f1d1b7696b9edd22","secret_hash":"0x40ddb2475d5b7a74ba9f28f6701e355c570f4e3d3fb965f352c67c17edd5e2b3","sc_user_address":"0x201b20123b3c489b47fde27ce5b451a0fa55fd60"}
digest=7f6b833a6e95bb5b84b860ca9c0b4a88825d6fe43dcdc3e68acbc0049f5b5861MCPrepareLockout SUCCESS
{
	"name": "Response",
	"request_id": "0x05a4cdd2d917ec48dce718ddc4b05425eeb8290ca9a9023e0c786a79dd8fdad1",
	"error_code": "0000",
	"error_msg": "success",
	"data": {
		"mc_chain_name": "ethereum",
		"secret_hash": "0x40ddb2475d5b7a74ba9f28f6701e355c570f4e3d3fb965f352c67c17edd5e2b3",
		"secret": "0x0000000000000000000000000000000000000000000000000000000000000000",
		"mc_user_address_hex": "0x201B20123b3C489b47Fde27ce5b451a0fA55FD60",
		"sc_user_address": "0x201b20123b3c489b47fde27ce5b451a0fa55fd60",
		"sc_token_address": "0x7b6cfc6a9edd6a1f81af4f8b860eebc26bf1ae4b",
		"amount": 20,
		"mc_expiration": 9237632,
		"sc_expiration": 3674731,
		"mc_lock_status": 1,
		"sc_lock_status": 1,
		"data": null,
		"notary_id_in_charge": 2,
		"start_time": 1607682547,
		"start_sc_block_number": 3673733,
		"end_time": 0,
		"end_sc_block_number": 0,
		"btc_prepare_lockout_tx_hash_hex": "",
		"btc_prepare_lockout_vout": 0,
		"btc_lock_script_hex": "",
		"cross_fee": 0
	}
}
```
其中req是向公证人提交的请求数据.

### 3.3 主链lockout
此步骤将侧链中的eth token再跨回主链,变成eth.

```
 ./dnc lo
```

```
start to lockout :
 ======> [chain=ethereum ]
[MC] lasted block number = 9237141
start to lockout :
 ======> [account=0x201b20123b3c489b47fde27ce5b451a0fa55fd60 secret=0x8bb14feefdcdf07fdbda4297be7e64f4494db5af895f9f0eceb39616bdf5c2b3 secretHash=0x40ddb2475d5b7a74ba9f28f6701e355c570f4e3d3fb965f352c67c17edd5e2b3]
lockout tx=
	TX(1fbc44ab3a86960eb6bed3b00389378c54c363cdfcfa36e58b890504d123c319)
	Contract: false
	From:     201b20123b3c489b47fde27ce5b451a0fa55fd60
	To:       7b6cfc6a9edd6a1f81af4f8b860eebc26bf1ae4b
	Nonce:    8
	GasPrice: 0x11520a2b8
	GasLimit  0xc6da
	Value:    0x0
	Data:     0x043d9180000000000000000000000000201b20123b3c489b47fde27ce5b451a0fa55fd608bb14feefdcdf07fdbda4297be7e64f4494db5af895f9f0eceb39616bdf5c2b3
	V:        0x29
	R:        0x786d940f49398fd1c8cb6f2d86718b493faff24df44d85fd9dbc46ce0b4ba91f
	S:        0x7c05adb1ab88af43bba8fefdab54a4e4b28939ced976b92f5b325f800f751699
	Hex:      f8a90885011520a2b882c6da947b6cfc6a9edd6a1f81af4f8b860eebc26bf1ae4b80b844043d9180000000000000000000000000201b20123b3c489b47fde27ce5b451a0fa55fd608bb14feefdcdf07fdbda4297be7e64f4494db5af895f9f0eceb39616bdf5c2b329a0786d940f49398fd1c8cb6f2d86718b493faff24df44d85fd9dbc46ce0b4ba91fa07c05adb1ab88af43bba8fefdab54a4e4b28939ced976b92f5b325f800f751699

Lockout SUCCESS
```
通过调用合约的lockout,将密码`0x78f6c1c7dd75b8e34bdce6496a368be86ab556c3579e49a9e532cf534abbef5c`公布在链上,这样公证人就可以在侧链上销毁相应的token.
该Tx为[1fbc44ab3a86960eb6bed3b00389378c54c363cdfcfa36e58b890504d123c319](https://ropsten.etherscan.io/tx/0x1fbc44ab3a86960eb6bed3b00389378c54c363cdfcfa36e58b890504d123c319)

#### 3.4 查询状态

```
./dnc q --all

===> MC/SC Lasted BlockNumber info :
[MC] lasted block number = 9237149
[SC] lasted block number = 3673745

===> MC/SC User account info :
[MC]user 201b2012 account balance : 997924827954760176
[SC]user 201b2012 sctoken balance : 0

===> MC/SC Contract account info :
[MC]contract 7b6cfc6a account balance : 0

===> MC/SC Contract data info :
[MC]data of lockin  : Empty
[SC]data of lockin  : Empty
[MC]data of lockout : Empty
[SC]data of lockout : Empty
```
可以看到主链合约不再锁定任何eth:`[MC]contract 7b6cfc6a account balance : 0`
侧链合约上账户`201b2012`也没有了任何代币.

