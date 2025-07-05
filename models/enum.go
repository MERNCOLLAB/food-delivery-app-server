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
	Pending          Status = "PENDING"
	AcceptedByOwner  Status = "ACCEPTED_BY_OWNER"
	AcceptedByDriver Status = "ACCEPTED_BY_DRIVER"
	RejectedByOwner  Status = "REJECTED_BY_OWNER"
	RejectedByDriver Status = "REJECTED_BY_DRIVER"
	ReadyForPickUp   Status = "READY_FOR_PICKUP"
	Assigned         Status = "ASSIGNED"
	InTransit        Status = "IN_TRANSIT"
	Delivered        Status = "DELIVERED"
	Canceled         Status = "ORDER_CANCELED"
)

type PaymentStatus string

const (
	Waiting PaymentStatus = "WAITING"
	Success PaymentStatus = "SUCCESS"
	Failed  PaymentStatus = "FAILED"
)
