package main

import "github.com/hanzhuoxian/mall/internal/userserver"

func main() {
	userserver.NewApp("mall-user-server").Run()
}
