package main

import (
	"fmt"

	"github.com/hanzhuoxian/mall/pkg/version"
)

func main() {
	fmt.Println("hello mall version: ", version.Get().ToJson())
}
