package enums

import "testing"

func TestMakeStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   MakeStatus
		expected int
	}{
		{"MakeStatusWaitMake", MakeStatusWaitMake, 0},
		{"MakeStatusMaking", MakeStatusMaking, 1},
		{"MakeStatusMade", MakeStatusMade, 2},
		{"MakeStatusMakeFail", MakeStatusMakeFail, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.status) != tt.expected {
				t.Errorf("Expected %s to be %d, but got %d", tt.name, tt.expected, int(tt.status))
			}
		})
	}
}

func TestGetMakeStatusDesc(t *testing.T) {
	tests := []struct {
		name     string
		status   MakeStatus
		expected string
	}{
		{"Wait make status", MakeStatusWaitMake, "待制作"},
		{"Making status", MakeStatusMaking, "制作中"},
		{"Made status", MakeStatusMade, "制作完成"},
		{"Make fail status", MakeStatusMakeFail, "制作失败"},
		{"Unknown status", MakeStatus(999), "未知状态"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMakeStatusDesc(tt.status)
			if result != tt.expected {
				t.Errorf("Expected GetMakeStatusDesc(%d) to be '%s', but got '%s'", tt.status, tt.expected, result)
			}
		})
	}
}

func TestMakeStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		status   MakeStatus
		expected string
	}{
		{"Wait make status string", MakeStatusWaitMake, "待制作"},
		{"Making status string", MakeStatusMaking, "制作中"},
		{"Made status string", MakeStatusMade, "制作完成"},
		{"Make fail status string", MakeStatusMakeFail, "制作失败"},
		{"Unknown status string", MakeStatus(999), "未知状态"},
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

func TestMakeStatus_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   MakeStatus
		expected bool
	}{
		{"Valid WaitMake status", MakeStatusWaitMake, true},
		{"Valid Making status", MakeStatusMaking, true},
		{"Valid Made status", MakeStatusMade, true},
		{"Valid MakeFail status", MakeStatusMakeFail, true},
		{"Invalid status -1", MakeStatus(-1), false},
		{"Invalid status 4", MakeStatus(4), false},
		{"Invalid large status", MakeStatus(999), false},
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
