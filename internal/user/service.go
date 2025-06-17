package user

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) UpdateUser() {

}

func (s *Service) UpdateProfilePicture() {

}

func (s *Service) DeleteUser() {

}

func (s *Service) GetAllUsers() {

}
