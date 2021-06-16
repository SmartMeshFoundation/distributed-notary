#!/bin/sh
abigen --sol sidechain_erc20.sol --pkg contracts --out ../chain/spectrum/contracts/SideChainErc20Token.go
abigen --sol mainchain_erc20.sol --pkg contracts --out ../chain/ethereum/contracts/MainChainErc20.go

