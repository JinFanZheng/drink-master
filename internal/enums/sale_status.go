package enums

// SaleStatus represents the sale status of a material silo
type SaleStatus int

const (
	// SaleStatusOff represents a material silo that is not available for sale
	SaleStatusOff SaleStatus = 0 // 停售
	// SaleStatusOn represents a material silo that is available for sale
	SaleStatusOn SaleStatus = 1 // 在售
)

// GetSaleStatusDesc returns the description of the sale status
func GetSaleStatusDesc(status SaleStatus) string {
	switch status {
	case SaleStatusOn:
		return "在售"
	case SaleStatusOff:
		return "停售"
	default:
		return "未知状态"
	}
}

// String returns the string representation of the sale status
func (ss SaleStatus) String() string {
	return GetSaleStatusDesc(ss)
}

// IsValid checks if the sale status is valid
func (ss SaleStatus) IsValid() bool {
	return ss >= SaleStatusOff && ss <= SaleStatusOn
}

// ToAPIString converts the sale status to API string representation
func (ss SaleStatus) ToAPIString() string {
	switch ss {
	case SaleStatusOn:
		return "On"
	case SaleStatusOff:
		return "Off"
	default:
		return "Off"
	}
}

// SaleStatusFromAPIString converts API string to SaleStatus enum
func SaleStatusFromAPIString(status string) SaleStatus {
	switch status {
	case "On":
		return SaleStatusOn
	case "Off":
		return SaleStatusOff
	default:
		return SaleStatusOff // default to off
	}
}
