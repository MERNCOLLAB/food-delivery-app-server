package menuitem

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateMenuItem() {

}

func (s *Service) GetMenuItemByRestaurant() {

}

func (s *Service) UpdateMenuItem() {
}

func (s *Service) DeleteMenuItem() {

}
