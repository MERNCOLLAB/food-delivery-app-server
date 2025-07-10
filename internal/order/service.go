package order

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) PlaceOrder(restaurantID string, orderReq PlaceOrderRequest) (*PlaceOrderResponse, error) {
	return nil, nil
}

func (s *Service) GetOrderByRestaurant() {

}

func (s *Service) GetOrderDetails() {

}

func (s *Service) GetOrderHistory() {

}

func (s *Service) GetAvailableOrders() {

}
