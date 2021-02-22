package repo

import (
	"fmt"
	"learnGo/wire/config"
)

// Repo ...
type Repo interface {
	GetHello(name string) string
}

// NewRepo ...
func NewRepo(cfg *config.Config) Repo {
	return &DefultRepo{
		Prefix: cfg.DBCfg.Prefix,
	}
}

// DefultRepo ...
type DefultRepo struct {
	Prefix string
}

// GetHello ~~
func (repo *DefultRepo) GetHello(name string) string {
	return fmt.Sprintf("%s %s", repo.Prefix, name)
}
