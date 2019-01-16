#!/bin/sh

#有可能上一个节点还在运行着没有退出 kill它
ps -ef | grep geth  |grep 7888| grep -v grep | awk '{print $2}' |xargs kill -9
ps -ef | grep geth  |grep 9888| grep -v grep | awk '{print $2}' |xargs kill -9

## 准备搭建私链
geth version
rm -rf 7888/geth
rm -rf 9888/geth
geth --datadir 7888 init poatestnet7888.json
geth --datadir 9888 init poatestnet9888.json

# 尽量避免不必要的log输出,干扰photon信息
geth --datadir=./7888 --unlock 3de45febbd988b6e417e4ebd2c69e42630fefbf0 --password ./7888/keystore/pass --port 7888 --networkid 7888 --ws --wsaddr 0.0.0.0 --wsorigins "*" --wsport 27888 --rpc --rpccorsdomain "*" --rpcapi eth,admin,web3,net,debug,personal --rpcport 17888 --rpcaddr 127.0.0.1 --mine  --verbosity 1 --nodiscover &

geth --datadir=./9888 --unlock 3de45febbd988b6e417e4ebd2c69e42630fefbf0 --password ./9888/keystore/pass --port 9888 --networkid 9888 --ws --wsaddr 0.0.0.0 --wsorigins "*" --wsport 29888 --rpc --rpccorsdomain "*" --rpcapi eth,admin,web3,net,debug,personal --rpcport 19888 --rpcaddr 127.0.0.1 --mine  --verbosity 1 --nodiscover &
