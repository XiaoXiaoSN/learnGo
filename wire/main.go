package main

import (
	// "learnGo/wire/repo"
	// "learnGo/wire/service"

	_ "github.com/google/wire"
)

func main() {
	// 原本的寫法要這樣
	// repo := repo.NewDefaultRepo()
	// svc := service.NewDefaultService(repo)
	// svc.SayHello("main package")

	// 現在寫一個 initService 來自動注入
	svc, err := initService()
	if err != nil {
		panic(err)
	}
	svc.SayHello("wire main package")
}
