// mall-user-server 是用户服务的可执行入口，负责启动用户服务器进程。
package main

import "github.com/hanzhuoxian/mall/internal/userserver"

// main 是程序入口，创建并运行用户服务应用。
func main() {
	userserver.NewApp("mall-user-server").Run()
}
