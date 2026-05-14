package main

import (
	"demo/internal/app"
	"os"
)

// @title			Demo Service
// @version			1.0
// @description		Демонстрация структуры сервиса и походов к разработке

// @host
// @BasePath /demo-service

// @tag.name admin
// @tag.description Админ API

// @tag.name api
// @tag.description Системное API

// @tag.name public
// @tag.description Публичное API

func main() {
	args := os.Args

	var cmd string
	if len(args) >= 2 {
		cmd = args[1]
	}

	switch cmd {
	default:
		app.NewApp().Run()
	}
}
