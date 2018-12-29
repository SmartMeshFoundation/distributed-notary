#!/bin/bash
dnotary --address=0x33df901abc22dcb7f33c2a77ad43cc98fbfa0790 --user-listen=127.0.0.1:3331 --notary-listen=127.0.0.1:33301 --notary-config-file=../notary.conf --keystore-path=../../keystore
