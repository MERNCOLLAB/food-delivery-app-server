package restaurant

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateRestaurant() {

}

func (s *Service) GetRestaurantByOwner() {

}

func (s *Service) UpdateRestaurant() {

}

func (s *Service) DeleteRestaurant() {

}
