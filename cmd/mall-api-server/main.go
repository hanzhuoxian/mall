package main

import "github.com/hanzhuoxian/mall/internal/apiserver"

func main() {
	apiserver.NewApp("mall-api-server").Run()
}
