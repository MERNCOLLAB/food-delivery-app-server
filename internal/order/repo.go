package order

import (
	"food-delivery-app-server/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateOrder(order *models.Order, items []models.OrderItem) (*models.Order, error) {
	// Starts a transaction to ensure order and order items creation are safe. If any one of those fails, the entire transaction
	// rolls back to maintain data integrity
	tx := r.db.Begin()
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for i := range items {
		items[i].OrderID = order.ID
		if err := tx.Create(&items[i]).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	return order, nil
}

func (r *Repository) GetMenuItemDetails(menuItemID uuid.UUID) (*models.MenuItem, error) {
	var menuItem models.MenuItem
	if err := r.db.Where("id = ?", menuItemID).First(&menuItem).Error; err != nil {
		return nil, err
	}
	return &menuItem, nil
}

func (r *Repository) GetRestaurantByID(restoId uuid.UUID) (*models.Restaurant, error) {
	var resto models.Restaurant
	if err := r.db.Where("id = ?", restoId).First(&resto).Error; err != nil {
		return nil, err
	}
	return &resto, nil
}

func (r *Repository) CreateNotification(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *Repository) GetOrderDetailsByID(orID uuid.UUID) (*models.Order, error) {
	var order models.Order

	err := r.db.
		Preload("OrderItems.MenuItem").
		Preload("Restaurant").
		Preload("Customer").
		Preload("Driver").
		First(&order, "id = ?", orID).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *Repository) GetOrderByRestaurantID() {

}

func (r *Repository) UpdateOrderStatus(orderID uuid.UUID, status string) error {
	var order models.Order
	if err := r.db.First(&order, "id = ?", orderID).Error; err != nil {
		return err
	}

	order.Status = models.Status(status)
	return r.db.Save(&order).Error
}

func (r *Repository) GetOrderHistoryByUser(uID uuid.UUID, role string) ([]models.Order, error) {
	var orders []models.Order

	query := r.db.
		Preload("OrderItems.MenuItem").
		Preload("Restaurant").
		Order("placed_at DESC")

	switch role {
	case "CUSTOMER":
		query = query.Where("customer_id = ?", uID)
	case "DRIVER":
		query = query.Where("driver_id = ?", uID)
	default:
		return nil, nil
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *Repository) GetUserRoleByID(uID uuid.UUID) (string, error) {
	var user models.User
	if err := r.db.Select("role").First(&user, "id = ?", uID).Error; err != nil {
		return "", err
	}
	return string(user.Role), nil
}

func (r *Repository) GetOrdersByRestaurantID(restoID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.
		Preload("OrderItems.MenuItem").
		Preload("Customer").
		Preload("Driver").
		Where("restaurant_id = ?", restoID).
		Order("placed_at DESC").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *Repository) GetOrderByStatus(status models.Status) ([]models.Order, error) {
	var orders []models.Order
	query := r.db.
		Preload("OrderItems.MenuItem").
		Preload("Restaurant").
		Preload("Customer").
		Preload("Driver").
		Where("status = ?", status)

	err := query.Order("placed_at DESC").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *Repository) GetOrderByStatusAndDriver(dr uuid.UUID, statuses []models.Status) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.
		Preload("OrderItems.MenuItem").
		Preload("Restaurant").
		Preload("Customer").
		Preload("Driver").
		Where("status IN ?", statuses).
		Where("driver_id = ?", dr).
		Order("placed_at DESC").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *Repository) UpdateOrderStatusAndDriver(order *models.Order) error {
	return r.db.Save(order).Error
}
