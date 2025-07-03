package address

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAddress() {

}

func (s *Service) GetAddress() {

}

func (s *Service) UpdateAddress() {

}

func (s *Service) DeleteAddress() {

}
