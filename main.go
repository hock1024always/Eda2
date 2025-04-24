package main

import "circuit_optimizer/router"

// 入口文件
func main() {
	r := router.Router()
	r.Run(":9090")
}
