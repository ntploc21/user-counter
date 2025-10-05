package main

import (
	"locntp-user-counter/config"
	"locntp-user-counter/internal/app"
)

func main() {
	config.Load()

	//appCfg := config.GetAppConfig()

	app.StartServer()
}
