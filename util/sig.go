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

// VerifySignature verifies an Ethereum signature given the message hash and the signature.
func VerifySignature(messageHash []byte, signature []byte) (bool, error) {
	// Recover the public key from the signature
	publicKeyECDSA, err := crypto.SigToPub(messageHash, signature)
	if err != nil {
		return false, err
	}
	recoveredPubKey := crypto.FromECDSAPub(publicKeyECDSA)

	return crypto.VerifySignature(recoveredPubKey, messageHash, signature[:64]), nil
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
