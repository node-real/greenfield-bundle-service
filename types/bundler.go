package types

import (
	"crypto/sha256"
	"fmt"
	"math/big"
)

// PickBundlerIndexForAccount deterministically selects a bundler index for a given account.
func PickBundlerIndexForAccount(bundlerCount int, account string) (int, error) {
	if bundlerCount <= 0 {
		return -1, fmt.Errorf("bundler count must be positive")
	}

	// Hash the account string
	hash := sha256.Sum256([]byte(account))

	// Convert the hash to a big integer and then mod it by the number of bundlers
	// This will give a consistent index based on the account string
	hashInt := new(big.Int).SetBytes(hash[:])
	index := new(big.Int).Mod(hashInt, big.NewInt(int64(bundlerCount)))

	return int(index.Int64()), nil
}
