package main

import (
	"fmt"

	"github.com/moevm/nosql1h25-writer/backend/pkg/hasher"
)

func main() {
	h := hasher.NewBcrypt()
	fmt.Println(h.Hash("password123"))
}
