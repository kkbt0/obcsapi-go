package main

import (
	"fmt"
	"time"
)

func main() {
	diff := time.Until(time.Now().AddDate(0, 0, 2))
	fmt.Println(diff)
	fmt.Println(int(diff.Hours() / 24))
}
