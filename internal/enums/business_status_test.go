package enums

import "testing"

func TestBusinessStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   BusinessStatus
		expected int
	}{
		{"BusinessStatusOpen", BusinessStatusOpen, 1},
		{"BusinessStatusClose", BusinessStatusClose, 2},
		{"BusinessStatusOffline", BusinessStatusOffline, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.status) != tt.expected {
				t.Errorf("Expected %s to be %d, but got %d", tt.name, tt.expected, int(tt.status))
			}
		})
	}
}

func TestGetBusinessStatusDesc(t *testing.T) {
	tests := []struct {
		name     string
		status   BusinessStatus
		expected string
	}{
		{"Open status", BusinessStatusOpen, "营业中"},
		{"Close status", BusinessStatusClose, "暂停营业"},
		{"Offline status", BusinessStatusOffline, "设备离线"},
		{"Invalid status", BusinessStatus(999), "未知状态"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetBusinessStatusDesc(tt.status)
			if result != tt.expected {
				t.Errorf("Expected GetBusinessStatusDesc(%d) to be '%s', but got '%s'", tt.status, tt.expected, result)
			}
		})
	}
}

func TestBusinessStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		status   BusinessStatus
		expected string
	}{
		{"Open status string", BusinessStatusOpen, "营业中"},
		{"Close status string", BusinessStatusClose, "暂停营业"},
		{"Offline status string", BusinessStatusOffline, "设备离线"},
		{"Invalid status string", BusinessStatus(999), "未知状态"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.String()
			if result != tt.expected {
				t.Errorf("Expected %d.String() to be '%s', but got '%s'", tt.status, tt.expected, result)
			}
		})
	}
}

func TestBusinessStatus_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   BusinessStatus
		expected bool
	}{
		{"Valid Open status", BusinessStatusOpen, true},
		{"Valid Close status", BusinessStatusClose, true},
		{"Valid Offline status", BusinessStatusOffline, true},
		{"Invalid status 0", BusinessStatus(0), false},
		{"Invalid status 4", BusinessStatus(4), false},
		{"Invalid negative status", BusinessStatus(-1), false},
		{"Invalid large status", BusinessStatus(999), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.IsValid()
			if result != tt.expected {
				t.Errorf("Expected %d.IsValid() to be %t, but got %t", tt.status, tt.expected, result)
			}
		})
	}
}
