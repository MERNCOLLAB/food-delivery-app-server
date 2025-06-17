package auth

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SignUp() {
}

func (s *Service) SignIn() {
}

func (s *Service) OAuth() {
}

func (s *Service) SignOut() {
}
