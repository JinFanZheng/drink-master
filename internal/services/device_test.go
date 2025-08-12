package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeviceService_CheckDeviceOnline(t *testing.T) {
	service := NewDeviceService()

	// 测试有效设备ID
	online, err := service.CheckDeviceOnline("device-123")
	require.NoError(t, err)
	assert.True(t, online) // 当前实现总是返回true

	// 测试空设备ID
	online, err = service.CheckDeviceOnline("")
	require.NoError(t, err)
	assert.False(t, online)
}

func TestDeviceService_UpdateRegister(t *testing.T) {
	service := NewDeviceService()

	params := map[string]int{
		"temperature": 25,
		"pressure":    100,
	}

	// 测试更新寄存器
	err := service.UpdateRegister("device-123", params)
	require.NoError(t, err)

	// 测试空参数
	err = service.UpdateRegister("device-123", nil)
	require.NoError(t, err)

	// 测试空设备ID
	err = service.UpdateRegister("", params)
	require.NoError(t, err)
}

func TestDeviceService_GetDeviceStatus(t *testing.T) {
	service := NewDeviceService()

	// 测试在线设备状态
	status, err := service.GetDeviceStatus("device-123")
	require.NoError(t, err)
	require.NotNil(t, status)
	assert.Equal(t, "device-123", status.DeviceID)
	assert.True(t, status.Online)
	assert.NotEmpty(t, status.LastSeen)

	// 验证LastSeen时间格式
	_, parseErr := time.Parse(time.RFC3339, status.LastSeen)
	assert.NoError(t, parseErr)

	// 测试离线设备（空ID）
	status, err = service.GetDeviceStatus("")
	require.NoError(t, err)
	require.NotNil(t, status)
	assert.Equal(t, "", status.DeviceID)
	assert.False(t, status.Online)
	assert.Empty(t, status.LastSeen)
}

func TestNewDeviceService(t *testing.T) {
	service := NewDeviceService()
	assert.NotNil(t, service)

}
