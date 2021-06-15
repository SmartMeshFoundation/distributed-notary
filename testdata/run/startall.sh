#!/usr/bin/env bash
# 删除历史数据
#rm -rf .notary*
#rm -f .n*.txt
# 启动两条测试链
#cd ../deploygeth
#sh ./deploygeth.sh
#cd ../run
#kill  历史dnotary
cd ../../cmd/dnotary 
go install
cd -
sh ./stopall.sh
ps -ef | grep dnotary  | grep -v grep | awk '{print $2}' |xargs kill -9
#MainEndPoint="http://193.112.248.133:19888"
#SideEndPoint="http://193.112.248.133:17888"
MainEndPoint="http://106.52.171.12:18003"
SideEndPoint="http://106.52.171.12:17004"
nonce_server -smc-rpc-endpoint=$SideEndPoint --eth-rpc-endpoint=$MainEndPoint  >>.nonce-server.txt 2>&1 &
#0
dnotary --address=0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa --user-listen=0.0.0.0:8030 --notary-listen=0.0.0.0:33300 --notary-config-file=./notary.conf --keystore-path=../keystore --datadir=./.notary_n0  --smc-rpc-endpoint=$SideEndPoint --eth-rpc-endpoint=$MainEndPoint --jettrade-eth-address=0x0E291c99d7A67cF64c03313efFc18Ca07eEffa1f --jettrade-spectrum-address=0xeCad09712c4b9f0c9c2356604ac3524d164C7f6B  >>.n0.txt 2>&1 &
#1
dnotary --address=0x33df901abc22dcb7f33c2a77ad43cc98fbfa0790 --user-listen=0.0.0.0:8031 --notary-listen=0.0.0.0:33301 --notary-config-file=./notary.conf --keystore-path=../keystore --datadir=./.notary_n1  --smc-rpc-endpoint=$SideEndPoint  --eth-rpc-endpoint=$MainEndPoint --jettrade-eth-address=0x0E291c99d7A67cF64c03313efFc18Ca07eEffa1f --jettrade-spectrum-address=0xeCad09712c4b9f0c9c2356604ac3524d164C7f6B >>.n1.txt 2>&1 &
#2
dnotary --address=0x8c1b2e9e838e2bf510ec7ff49cc607b718ce8401 --user-listen=0.0.0.0:8032 --notary-listen=0.0.0.0:33302 --notary-config-file=./notary.conf --keystore-path=../keystore --datadir=./.notary_n2  --smc-rpc-endpoint=$SideEndPoint   --eth-rpc-endpoint=$MainEndPoint --jettrade-eth-address=0x0E291c99d7A67cF64c03313efFc18Ca07eEffa1f --jettrade-spectrum-address=0xeCad09712c4b9f0c9c2356604ac3524d164C7f6B >>.n2.txt 2>&1 &
#3
dnotary --address=0xc4c08f9227be0f1750f5d5467eed462ec133b15e --user-listen=0.0.0.0:8033 --notary-listen=0.0.0.0:33303 --notary-config-file=./notary.conf --keystore-path=../keystore --datadir=./.notary_n3  --smc-rpc-endpoint=$SideEndPoint  --eth-rpc-endpoint=$MainEndPoint --jettrade-eth-address=0x0E291c99d7A67cF64c03313efFc18Ca07eEffa1f --jettrade-spectrum-address=0xeCad09712c4b9f0c9c2356604ac3524d164C7f6B >>.n3.txt 2>&1 &
#4
dnotary --address=0x543fc024cdd1f0d346a306f5e99ec0d8fe392920 --user-listen=0.0.0.0:8034 --notary-listen=0.0.0.0:33304 --notary-config-file=./notary.conf --keystore-path=../keystore --datadir=./.notary_n4  --smc-rpc-endpoint=$SideEndPoint  --eth-rpc-endpoint=$MainEndPoint --jettrade-eth-address=0x0E291c99d7A67cF64c03313efFc18Ca07eEffa1f --jettrade-spectrum-address=0xeCad09712c4b9f0c9c2356604ac3524d164C7f6B >>.n4.txt 2>&1 &
#5
dnotary --address=0x920a90acc9164272ede4ae1e9c33841f019f53a4 --user-listen=0.0.0.0:8035 --notary-listen=0.0.0.0:33305 --notary-config-file=./notary.conf --keystore-path=../keystore --datadir=./.notary_n5  --smc-rpc-endpoint=$SideEndPoint  --eth-rpc-endpoint=$MainEndPoint --jettrade-eth-address=0x0E291c99d7A67cF64c03313efFc18Ca07eEffa1f --jettrade-spectrum-address=0xeCad09712c4b9f0c9c2356604ac3524d164C7f6B >>.n5.txt 2>&1 &
#6
dnotary --address=0x215c0d259ac31571a43295f2e411a697cd30748c --user-listen=0.0.0.0:8036 --notary-listen=0.0.0.0:33306 --notary-config-file=./notary.conf --keystore-path=../keystore --datadir=./.notary_n6  --smc-rpc-endpoint=$SideEndPoint  --eth-rpc-endpoint=$MainEndPoint --jettrade-eth-address=0x0E291c99d7A67cF64c03313efFc18Ca07eEffa1f --jettrade-spectrum-address=0xeCad09712c4b9f0c9c2356604ac3524d164C7f6B >>.n6.txt 2>&1 &



