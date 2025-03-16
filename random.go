package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func generateNewCommitment() (string, string) {
	secretBuff := _rand_Uint256()
	commitment := _hashCommitment(secretBuff)
	encSecret := hex.EncodeToString(secretBuff)
	return encSecret, commitment
}

// Random Uint256 buffer
func _rand_Uint256() []byte {
	buff := make([]byte, 32)
	rand.Read(buff)
	return buff
}

// Sha256 commitment
func _hashCommitment(buff []byte) string {
	hash := sha256.Sum256(buff)
	hexHash := hex.EncodeToString(hash[:])
	return hexHash
}
