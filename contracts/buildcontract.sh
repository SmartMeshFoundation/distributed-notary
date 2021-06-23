#!/bin/sh
abigen --sol sidechain_erc20.sol --pkg contracts --out ../chain/heco/contracts/SideChainErc20Token.go
abigen --sol mainchain_erc20.sol --pkg contracts --out ../chain/spectrum/contracts/MainChainErc20.go
