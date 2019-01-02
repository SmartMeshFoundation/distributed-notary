#!/bin/bash
dnotary --address=0xc4c08f9227be0f1750f5d5467eed462ec133b15e --user-listen=127.0.0.1:3333 --notary-listen=127.0.0.1:33303 --notary-config-file=../notary.conf --keystore-path=../../keystore
