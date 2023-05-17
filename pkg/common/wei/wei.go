package wei

// WeiToEth converts wei to eth
func WeiToEth(wei uint64) float64 {
	return float64(wei) / 1000000000000000000
}

// EthToWei converts eth to wei
func EthToWei(eth float64) uint64 {
	return uint64(eth * 1000000000000000000)
}

// GweiToWei converts gwei to wei
func GweiToWei(gwei uint64) uint64 {
	return gwei * 1000000000
}

// WeiToGwei converts wei to gwei
func WeiToGwei(wei uint64) uint64 {
	return wei / 1000000000
}
