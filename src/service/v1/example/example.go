package example

import "github.com/forkyid/go-consumer-boilerplate/src/repository/v1/example"

type Service struct {
	repo example.Repositorier
}

func NewService(
	repo example.Repositorier,
) *Service {
	return &Service{
		repo: repo,
	}
}

type Servicer interface {
}
