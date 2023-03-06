package main

import (
	"fmt"
	"testing"
)

func TestRomanTokne(t *testing.T) {
	token := GengerateToken(99)
	fmt.Println(len(token))
	fmt.Println(token)
}

func TestModToken(t *testing.T) {
	token1 := Token{TokenString: GengerateToken(10), GenerateTime: GengerateToken(20)}
	ModTokenFile(token1, "token1")
	token2, _ := GetToken("token1")
	fmt.Println(token1, token2)
}
