package auth

import "fmt"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) TestError() error {
	return fmt.Errorf("failed to sign up, test error")
}

func (s *Service) SignUp() {
}

func (s *Service) SignIn() {
}

func (s *Service) OAuth() {
}

func (s *Service) SignOut() {
}
