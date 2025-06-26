package resetpassword

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RequestResetPassword() {

}

func (s *Service) VerifyResetCode() {

}

func (s *Service) UpdatePassword() {

}
