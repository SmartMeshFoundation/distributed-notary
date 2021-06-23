package cfg

// GetCfgByChainName 方便使用
func GetCfgByChainName(name string) *ChainCfg {
	if name == SMC.Name {
		return SMC
	} else if name == HECO.Name {
		return HECO
	}
	panic("wrong code")
}

// GetCrossFeeRate :
func GetCrossFeeRate(chainName string) int64 {
	c := GetCfgByChainName(chainName)
	return c.CrossFeeRate
}

// GetMinExpirationBlock4User :
func GetMinExpirationBlock4User(chainName string) uint64 {
	c := GetCfgByChainName(chainName)
	return uint64(Cross.MinExpirationTime4User / c.BlockPeriod)
}

// GetMinExpirationBlock4Notary :
func GetMinExpirationBlock4Notary(chainName string) uint64 {
	c := GetCfgByChainName(chainName)
	return uint64(Cross.MinExpirationTime4Notary / c.BlockPeriod)
}
