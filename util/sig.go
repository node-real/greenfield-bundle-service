package util

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// GenerateRandomAccount generates a new private key and returns the private key and address in byte format
func GenerateRandomAccount() ([]byte, common.Address, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, common.Address{}, err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	return privateKeyBytes, crypto.PubkeyToAddress(privateKey.PublicKey), nil
}

// RecoverAddress recovers the Ethereum address from the given message hash and signature
func RecoverAddress(messageHash common.Hash, signature []byte) (common.Address, error) {
	// Recover the public key from the signature
	publicKey, err := crypto.SigToPub(messageHash.Bytes(), signature)
	if err != nil {
		return common.Address{}, err
	}

	// Extract the address from the public key
	address := crypto.PubkeyToAddress(*publicKey)
	return address, nil
}
