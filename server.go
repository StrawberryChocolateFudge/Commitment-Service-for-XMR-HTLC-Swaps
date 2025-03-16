package main

import "fmt"

func main() {
	secret, commitment := generateNewCommitment()
	fmt.Println(secret)
	fmt.Println(commitment)
}
