package service

import (
	"fmt"
	"learnGo/wire/repo"
)

// Service ...
type Service interface {
	SayHello(name string)
}

// NewService ...
func NewService(repo repo.Repo) Service {
	return &DefultService{
		repo: repo,
	}
}

// DefultService ...
type DefultService struct {
	repo repo.Repo
}

// SayHello ~~
func (svc *DefultService) SayHello(name string) {
	fmt.Println(svc.repo.GetHello(name))
}
