
# dnc instruction

## 1.  Initialization
### 1.1 dnc Usage help
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
### 1.2 dnc Configuration
Please use the attached dnc.json file, and ignore any content about btc in it.

### 1.3 Get the sidechain token contract address
```
./dnc config refresh
```
At this time, open dnc.json and you can see:
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
sc_token`0x7b6cfc6a9edd6a1f81af4f8b860eebc26bf1ae4b` in sc_token_list is the contract address which we want to use.

#### 1.4 Switch account
Please place the corresponding file in the keystore, and modify the `smc_user_address` and `eth_user_address` and the corresponding password in section 1.3.
Do not use the `dnc config reset` command to reset the configuration file.

#### 1.5 Query the current status of the contract

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
## 2. Main chain to side chain cross-chain

### 2.1 Main chain PrepareLockIn

```
./dnc pli --amount 20
```
The amount here is 20 wei,wei is the unit of measurement

After success, You will get the following output:
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
Parameters used in this cross-chain:
- secret:0xf1b820fd53c8747bd959f7172d5d38fe5f0174d1cfc15a7db3205b83584dadd2 ,This is private and will not be public on the chain at present
- secrethash: 0x474ce39fa753dbe30fe4d43e37a0959b2713a291a22053a88981b37ad7d90be1
- Cross-chain amount: 20
- Expiration: After 900 blocks 
- Tx: [e811c7719f51c1db85d0f53c3fbde0480db8940f8520e676a1b386eef44781d3](https://ropsten.etherscan.io/tx/0xe811c7719f51c1db85d0f53c3fbde0480db8940f8520e676a1b386eef44781d3)

### 2.2 Notify the notary

Notify the notary that I have just performed a cross-chain behavior and need a notary to preparelockin on the side chain
```
 ./dnc scpli
```

Get the following output:
```
req={"name":"User-SCPrepareLockin","request_id":"0xac97529f278d985a9be69cd6013cc52d038212a1ce703c24a536dd2480dc6c45","sc_token_address":"0x7b6cfc6a9edd6a1f81af4f8b860eebc26bf1ae4b","signer":"03f2307da626a34b0ebd1732b97df2517630bf699198c43b54f1d1b7696b9edd22","secret_hash":"0x474ce39fa753dbe30fe4d43e37a0959b2713a291a22053a88981b37ad7d90be1","mc_user_address":"0x201b20123b3c489b47fde27ce5b451a0fa55fd60"}
digest=427d86fc0c453fcae33a00ec7a23f9614a3fca94d1592c896a7111cd7e0d0c85
SCPrepareLockin SUCCESS
```
req is the data submitted to the notary. See the documentation for details.

### 2.3 Sidechain lockin to get Erc20 token
Announce the secret `0xf1b820fd53c8747bd959f7172d5d38fe5f0174d1cfc15a7db3205b83584dadd2` by calling lockin on the side chain
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
among them:
- Published secret: 0xf1b820fd53c8747bd959f7172d5d38fe5f0174d1cfc15a7db3205b83584dadd2
- Tx [7dc1f7b614f9432bc892d5243657a9853101b2e6735a22081f5695643a7abb0c](https://chain.smartmesh.cn/tx.html?hash=0x7dc1f7b614f9432bc892d5243657a9853101b2e6735a22081f5695643a7abb0c)

#### 2.4 query result
At this point, from the main chain to the side chain, you can see the result.
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

`[SC]user 201b2012 sctoken balance : 20` shows that account `201b2012` has 20wei tokens on the side chain
`[MC]contract 7b6cfc6a account balance : 20` indicates that 20wei eth is locked on the main chain contract.

## 3. Side chain to main chain cross-chain
The side chain crosses back to the main chain, which is to exchange the eth tokens just obtained back to eth.

### 3.1 Side chain preparelockout 
```
./dnc plo --amount 20
```

Output:
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
The parameters used in this cross-chain:
- secret: 0x78f6c1c7dd75b8e34bdce6496a368be86ab556c3579e49a9e532cf534abbef5c
- secrethash: 0xcede52a047a88b6c133ac79f83a606ce8e854a40af52eac4d153bbe315d0feb4
- Relative number of expired blocks: 1000
- amount: 20wei
- Tx [c758dfab62b4654abf11b45395ae84f93ca91487be1d202b4654a716cc123843](https://chain.smartmesh.cn/tx.html?hash=0xc758dfab62b4654abf11b45395ae84f93ca91487be1d202b4654a716cc123843)

### 3.2 Notify the notary

```
./dnc mcplo
```
Output:
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
 req is the requested data submitted to the notary.

### 3.3 main chain lockout
In this step, the eth token in the side chain is crossed back to the main chain to obtain eth.

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
By calling the lockout of the contract, the secret `0x78f6c1c7dd75b8e34bdce6496a368be86ab556c3579e49a9e532cf534abbef5c` is published on the chain, so that the notary can destroy the corresponding token on the side chain.
The Tx is [1fbc44ab3a86960eb6bed3b00389378c54c363cdfcfa36e58b890504d123c319](https://ropsten.etherscan.io/tx/0x1fbc44ab3a86960eb6bed3b00389378c54c363cdfcfa36e58b890504d123c319)

#### 3.4 Query status

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
You can see that the main chain contract no longer locks any eth:`[MC]contract 7b6cfc6a account balance : 0`
The account `201b2012` on the sidechain contract also does not have any tokens.

