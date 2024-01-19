package util

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// RecoverAddr recovers the sender address from msg and signature
func RecoverAddr(msg []byte, sig []byte) (common.Address, error) {
	pubKeyByte, err := secp256k1.RecoverPubkey(msg, sig)
	if err != nil {
		return common.Address{}, err
	}
	pubKey, _ := crypto.UnmarshalPubkey(pubKeyByte)
	address := crypto.PubkeyToAddress(*pubKey)
	return address, nil
}

func SignMessage(privateKeyBytes []byte, message []byte) ([]byte, error) {
	// Convert bytes to ECDSA private key
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	println("hex message", hex.EncodeToString(message))

	// Sign the message
	signature, err := crypto.Sign(message, privateKey)
	if err != nil {
		return nil, err
	}
	println("hex signature", hex.EncodeToString(signature))
	return signature, err
}

func TestVerifySignature(t *testing.T) {
	privateKey, ori, err := GenerateRandomAccount()
	if err != nil {
		t.Fatal(err)
	}
	println("priv: ", hex.EncodeToString(privateKey))
	println(ori.String())

	message := []byte("Hello World")
	messageHash := crypto.Keccak256(message)

	anotherMessage := []byte("Hello World, too")
	anotherMessageHash := crypto.Keccak256(anotherMessage)

	signature, err := SignMessage(privateKey, messageHash)
	if err != nil {
		t.Fatal(err)
	}

	if addr, err := RecoverAddr(anotherMessageHash, signature); err != nil {
		t.Fatal(err)
	} else {
		println(addr.String())
	}

	if addr, err := RecoverAddress(common.BytesToHash(anotherMessageHash), signature); err != nil {
		t.Fatal(err)
	} else {
		println(addr.String())
	}
}
