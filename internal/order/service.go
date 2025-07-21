package order

import (
	appErr "food-delivery-app-server/pkg/errors"
	"time"

	"food-delivery-app-server/models"
	"food-delivery-app-server/pkg/notifications"
	"food-delivery-app-server/pkg/utils"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) PlaceOrder(restaurantID string, userID string, orderReq PlaceOrderRequest) (*PlaceOrderResponse, error) {
	var totalAmount float64
	var orderItems []models.OrderItem

	restoID, err := utils.ParseId(restaurantID)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid Restaurant ID", err)
	}

	uID, err := utils.ParseId(userID)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid User ID", err)
	}

	for _, item := range orderReq.Items {
		menuItem, err := s.repo.GetMenuItemDetails(item.MenuItem)
		if err != nil {
			return nil, err
		}
		itemTotal := menuItem.Price * float64(item.Quantity)
		totalAmount += itemTotal

		orderItems = append(orderItems, models.OrderItem{
			ID:         utils.GenerateUUID(),
			MenuItemID: menuItem.ID,
			Quantity:   int32(item.Quantity),
			Price:      menuItem.Price,
		})
	}

	deliveryFee := 50.0
	totalAmount += deliveryFee

	orderID := utils.GenerateUUID()

	order := &models.Order{
		ID:              orderID,
		RestaurantID:    restoID,
		CustomerID:      &uID,
		Status:          models.Status("PENDING"),
		TotalAmount:     totalAmount,
		DeliveryFee:     deliveryFee,
		DeliveryAddress: orderReq.DeliveryAddress,
		PlacedAt:        time.Now(),
	}

	createdOrder, err := s.repo.CreateOrder(order, orderItems)
	if err != nil {
		return nil, err
	}

	resto, err := s.repo.GetRestaurantByID(restoID)
	if err == nil && resto.OwnerID != uuid.Nil {
		notification := &models.Notification{
			ID:        utils.GenerateUUID(),
			UserID:    resto.OwnerID,
			OrderID:   &createdOrder.ID,
			Message:   "You have a new order!",
			IsRead:    false,
			CreatedAt: time.Now(),
		}

		_ = s.repo.CreateNotification(notification)
	}

	orderRes := &PlaceOrderResponse{
		OrderID:         createdOrder.ID,
		Status:          createdOrder.Status,
		TotalAmount:     createdOrder.TotalAmount,
		DeliveryFee:     createdOrder.DeliveryFee,
		DeliveryAddress: createdOrder.DeliveryAddress,
		PlacedAt:        createdOrder.PlacedAt.Format(time.RFC1123),
		Items:           orderItems,
	}

	return orderRes, nil
}

func (s *Service) GetOrderByRestaurant(restaurantID string) ([]models.Order, error) {
	restoID, err := utils.ParseId(restaurantID)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid restaurant ID", err)
	}

	orders, err := s.repo.GetOrdersByRestaurantID(restoID)
	if err != nil {
		return nil, appErr.NewInternal("Failed to get orders for restaurant", err)
	}

	return orders, nil
}

func (s *Service) GetOrderDetails(orderId string) (*models.Order, error) {
	orID, err := utils.ParseId(orderId)
	if err != nil {
		return nil, appErr.NewInternal("Invalid ID", err)
	}

	order, err := s.repo.GetOrderDetailsByID(orID)
	if err != nil {
		return nil, appErr.NewInternal("Failed to query order details", err)
	}

	return order, nil
}

func (s *Service) GetOrderHistory(userId string) ([]models.Order, error) {
	uID, err := utils.ParseId(userId)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid user ID", err)
	}

	role, err := s.repo.GetUserRoleByID(uID)
	if err != nil {
		return nil, appErr.NewInternal("Failed to get user role", err)
	}

	orders, err := s.repo.GetOrderHistoryByUser(uID, role)
	if err != nil {
		return nil, appErr.NewInternal("Failed to get the order history", err)
	}

	return orders, nil
}

func (s *Service) GetAvailableOrders() ([]models.Order, error) {
	status := models.ReadyForPickUp
	orders, err := s.repo.GetOrderByStatus(status)
	if err != nil {
		return nil, appErr.NewInternal("Failed to get available orders", err)
	}
	return orders, nil
}

func (s *Service) GetAssignedOrders(driverId string) ([]models.Order, error) {
	dr, err := utils.ParseId(driverId)
	if err != nil {
		return nil, appErr.NewBadRequest("Invalid driver ID", err)
	}

	validStatuses := []models.Status{
		models.AcceptedByDriver,
		models.Assigned,
		models.InTransit,
		models.Delivered,
	}

	orders, err := s.repo.GetOrderByStatusAndDriver(dr, validStatuses)
	if err != nil {
		return nil, appErr.NewInternal("Failed to get assigned orders", err)
	}

	return orders, nil
}

func (s *Service) UpdateOrderStatus(req UpdateOrderStatusRequest, orderId string) error {
	orID, err := utils.ParseId(orderId)
	if err != nil {
		return appErr.NewBadRequest("Invalid order ID", err)
	}

	order, err := s.repo.GetOrderDetailsByID(orID)
	if err != nil {
		return appErr.NewInternal("Order not found", err)
	}

	currStatus := order.Status
	newStatus := models.Status(req.Status)

	allowedNext, ok := allowedStatusTransitions[currStatus]
	if !ok {
		return appErr.NewBadRequest("Invalid current order status", nil)
	}

	valid := false
	for _, s := range allowedNext {
		if s == newStatus {
			valid = true
			break
		}
	}

	if !valid {
		return appErr.NewBadRequest("Invalid status request", nil)
	}

	if err := s.repo.UpdateOrderStatus(orID, req.Status); err != nil {
		return appErr.NewInternal("Failed to update the order status", err)
	}

	_ = notifications.CreateStatusChangeNotifications(s.repo, order, newStatus)

	return nil
}

func (s *Service) UpdateDriverOrderStatus(req UpdateOrderStatusRequest, orderId string, driverId string) error {
	if req.Status != "ACCEPTED_BY_DRIVER" && req.Status != "REJECTED_BY_DRIVER" && req.Status != "IN_TRANSIT" && req.Status != "DELIVERED" {
		return appErr.NewBadRequest("Invalid request status by driver", nil)
	}

	orID, err := utils.ParseId(orderId)
	if err != nil {
		return appErr.NewBadRequest("Invalid order ID", err)
	}

	drID, err := utils.ParseId(driverId)
	if err != nil {
		return appErr.NewBadRequest("Invalid driver ID", err)
	}

	order, err := s.repo.GetOrderDetailsByID(orID)
	if err != nil {
		return appErr.NewInternal("Order not found", err)
	}

	switch req.Status {
	case string(models.AcceptedByDriver):
		if order.Status != models.ReadyForPickUp || order.DriverID != nil {
			return appErr.NewBadRequest("Order not available for driver assignment", nil)
		}
		order.Status = models.AcceptedByDriver
		order.DriverID = &drID
	case string(models.RejectedByDriver):
		if order.Status != models.ReadyForPickUp || order.DriverID != nil {
			return appErr.NewBadRequest("Order not available for driver assignment", nil)
		}
		order.Status = models.RejectedByDriver
	case string(models.InTransit):
		if order.Status != models.AcceptedByDriver || order.DriverID == nil || *order.DriverID != drID {
			return appErr.NewBadRequest("Order must be accepted by this driver before marking as in transit", nil)
		}
		order.Status = models.InTransit
	case string(models.Delivered):
		if order.Status != models.InTransit || order.DriverID == nil || *order.DriverID != drID {
			return appErr.NewBadRequest("Order must be in transit by this driver before marking as delivered", nil)
		}
		order.Status = models.Delivered
	default:
		return appErr.NewBadRequest("Invalid status transition", nil)
	}

	if req.Status == string(models.AcceptedByDriver) {
		order.Status = models.AcceptedByDriver
		order.DriverID = &drID
	} else if req.Status == string(models.RejectedByDriver) {
		order.Status = models.RejectedByDriver
	}

	if err := s.repo.UpdateOrderStatusAndDriver(order); err != nil {
		return appErr.NewInternal("Failed to update order status", err)
	}

	_ = notifications.CreateStatusChangeNotifications(s.repo, order, order.Status)

	return nil
}

func (s *Service) CancelOrder(orderId string, userId string) error {
	orID, err := utils.ParseId(orderId)
	if err != nil {
		return appErr.NewBadRequest("Invalid order ID", err)
	}

	uID, err := utils.ParseId(userId)
	if err != nil {
		return appErr.NewBadRequest("Invalid user ID", err)
	}

	order, err := s.repo.GetOrderDetailsByID(orID)
	if err != nil {
		return appErr.NewInternal("Order not found", err)
	}

	if order.Status != models.Pending && order.Status != models.AcceptedByOwner {
		return appErr.NewBadRequest("Order cannot be canceled at this stage", nil)
	}

	if order.CustomerID == nil || *order.CustomerID != uID {
		return appErr.NewUnauthorized("You are not allowed to cancel this order", nil)
	}

	if order.Status == models.AcceptedByOwner {
		if time.Since(order.UpdatedAt) > 3*time.Minute {
			return appErr.NewBadRequest("You can only cancel within 3 minutes after acceptance", nil)
		}
	}

	cancelStatus := string(models.Canceled)
	if err := s.repo.UpdateOrderStatus(orID, cancelStatus); err != nil {
		return appErr.NewInternal("Failed to cancel the order", err)
	}

	_ = notifications.CreateStatusChangeNotifications(s.repo, order, models.Canceled)

	return nil
}
