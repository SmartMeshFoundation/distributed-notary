#!/bin/bash
dnotary --address=0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa --user-listen=http://127.0.0.1:3330 --notary-listen=http://127.0.0.1:33300 --notary-config-file=../notary.conf --keystore-path=../../keystore
