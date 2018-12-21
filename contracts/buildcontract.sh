#!/bin/sh
abigen --sol sidechain_erc20.sol --pkg contracts --out ../spectrum/contracts/SideChainErc20Token.go
abigen --sol mainchain_erc20.sol --pkg contracts --out ../ethereum/contracts/MainChainErc20.go
