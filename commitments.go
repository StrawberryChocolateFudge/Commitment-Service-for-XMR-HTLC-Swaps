package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/iden3/go-iden3-crypto/poseidon"
)

func generateNewCommitment(hashfunc string) (string, string) {

	if hashfunc == "sha256" {
		secretBuff := _rand_Uint256()
		commitment := _hashSha256Commitment(secretBuff)
		encSecret := hex.EncodeToString(secretBuff)
		return encSecret, commitment
	}

	if hashfunc == "poseidon" {
		secretBuff := _rand_Uint254()
		commitment := _hashPoseidon(secretBuff)
		encSecret := hex.EncodeToString(secretBuff)
		return encSecret, commitment
	}

	return "", ""
}

// Random Uint256 buffer
func _rand_Uint256() []byte {
	buff := make([]byte, 32)
	rand.Read(buff)
	return buff
}

// A 31 byte buffer for the BN254 curve
func _rand_Uint254() []byte {
	buff := make([]byte, 31)
	rand.Read(buff)
	return buff
}

// Sha256 commitment
func _hashSha256Commitment(buff []byte) string {
	hash := sha256.Sum256(buff)
	return hex.EncodeToString(hash[:])
}

// Poseidon hash
func _hashPoseidon(buff []byte) string {
	hash := poseidon.Sum(buff)
	return hex.EncodeToString(hash[:])
}
