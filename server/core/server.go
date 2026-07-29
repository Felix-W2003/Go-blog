package core

import (
	"server/global"
	"server/initialize"
)

type server interface {
	ListenAndServe() error
}

// RunServer 用于启动服务器
func RunServer() {
	addr := global.Config.System.Addr()
	Router := initialize.InitRouter()

	// 加在所有JWT黑名单，存入本地缓存
	//
}
