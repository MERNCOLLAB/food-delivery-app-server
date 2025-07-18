package address

type CreateAddressRequest struct {
	Address   string `json:"address"`
	Label     string `json:"label"`
	IsDefault bool   `json:"isDefault"`
}

type UpdateAddressRequest struct {
	Address   *string `json:"address,omitempty"`
	Label     *string `json:"label,omitempty"`
	IsDefault *bool   `json:"isDefault,omitempty"`
}
