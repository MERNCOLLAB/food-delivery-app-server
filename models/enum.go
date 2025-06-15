package models

type Role string

const (
	Admin    Role = "ADMIN"
	Owner    Role = "OWNER"
	Driver   Role = "DRIVER"
	Customer Role = "CUSTOMER"
)

type Status string

const (
	Pending        Status = "PENDING"
	Accepted       Status = "ACCEPTED"
	Rejected       Status = "REJECTED"
	ReadyForPickUp Status = "READY_FOR_PICKUP"
	Assigned       Status = "ASSIGNED"
	InTransit      Status = "IN_TRANSIT"
	Delivered      Status = "DELIVERED"
	Canceled       Status = "CANCELED"
)

type PaymentStatus string

const (
	Waiting PaymentStatus = "WAITING"
	Success PaymentStatus = "SUCCESS"
	Failed  PaymentStatus = "FAILED"
)
