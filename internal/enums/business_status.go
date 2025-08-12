package enums

// BusinessStatus represents the business status of a machine
type BusinessStatus int

const (
	// BusinessStatusOpen represents a machine that is open for business
	BusinessStatusOpen BusinessStatus = 1 // 营业中
	// BusinessStatusClose represents a machine that is temporarily closed
	BusinessStatusClose BusinessStatus = 2 // 暂停营业
	// BusinessStatusOffline represents a machine that is offline
	BusinessStatusOffline BusinessStatus = 3 // 设备离线
)

// GetBusinessStatusDesc returns the description of the business status
func GetBusinessStatusDesc(status BusinessStatus) string {
	switch status {
	case BusinessStatusOpen:
		return "营业中"
	case BusinessStatusClose:
		return "暂停营业"
	case BusinessStatusOffline:
		return "设备离线"
	default:
		return "未知状态"
	}
}

// String returns the string representation of the business status
func (bs BusinessStatus) String() string {
	return GetBusinessStatusDesc(bs)
}

// IsValid checks if the business status is valid
func (bs BusinessStatus) IsValid() bool {
	return bs >= BusinessStatusOpen && bs <= BusinessStatusOffline
}
