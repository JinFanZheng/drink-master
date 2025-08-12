package enums

// MakeStatus represents the making status of a drink order
type MakeStatus int

const (
	// MakeStatusWaitMake represents an order waiting to be made
	MakeStatusWaitMake MakeStatus = 0 // 待制作
	// MakeStatusMaking represents an order currently being made
	MakeStatusMaking MakeStatus = 1 // 制作中
	// MakeStatusMade represents an order that has been completed
	MakeStatusMade MakeStatus = 2 // 制作完成
	// MakeStatusMakeFail represents an order that failed to be made
	MakeStatusMakeFail MakeStatus = 3 // 制作失败
)

// GetMakeStatusDesc returns the description of the make status
func GetMakeStatusDesc(status MakeStatus) string {
	switch status {
	case MakeStatusWaitMake:
		return "待制作"
	case MakeStatusMaking:
		return "制作中"
	case MakeStatusMade:
		return "制作完成"
	case MakeStatusMakeFail:
		return "制作失败"
	default:
		return "未知状态"
	}
}

// String returns the string representation of the make status
func (ms MakeStatus) String() string {
	return GetMakeStatusDesc(ms)
}

// IsValid checks if the make status is valid
func (ms MakeStatus) IsValid() bool {
	return ms >= MakeStatusWaitMake && ms <= MakeStatusMakeFail
}
