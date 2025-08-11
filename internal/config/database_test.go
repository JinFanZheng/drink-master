package config

import (
	"os"
	"testing"
)

func TestLoadDatabaseConfig(t *testing.T) {
	// 保存原始环境变量
	originalValues := map[string]string{
		"DB_HOST":     os.Getenv("DB_HOST"),
		"DB_PORT":     os.Getenv("DB_PORT"),
		"DB_USER":     os.Getenv("DB_USER"),
		"DB_PASSWORD": os.Getenv("DB_PASSWORD"),
		"DB_NAME":     os.Getenv("DB_NAME"),
	}

	defer func() {
		// 恢复原始环境变量
		for key, value := range originalValues {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	// 测试默认值
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")

	config := LoadDatabaseConfig()

	if config.Host != "localhost" {
		t.Errorf("Expected Host to be 'localhost', got '%s'", config.Host)
	}
	if config.Port != "3306" {
		t.Errorf("Expected Port to be '3306', got '%s'", config.Port)
	}
	if config.User != "drink_master" {
		t.Errorf("Expected User to be 'drink_master', got '%s'", config.User)
	}
	if config.Password != "" {
		t.Errorf("Expected Password to be empty, got '%s'", config.Password)
	}
	if config.DBName != "drink_master_dev" {
		t.Errorf("Expected DBName to be 'drink_master_dev', got '%s'", config.DBName)
	}

	// 测试自定义环境变量
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "3307")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")

	config = LoadDatabaseConfig()

	if config.Host != "testhost" {
		t.Errorf("Expected Host to be 'testhost', got '%s'", config.Host)
	}
	if config.Port != "3307" {
		t.Errorf("Expected Port to be '3307', got '%s'", config.Port)
	}
	if config.User != "testuser" {
		t.Errorf("Expected User to be 'testuser', got '%s'", config.User)
	}
	if config.Password != "testpass" {
		t.Errorf("Expected Password to be 'testpass', got '%s'", config.Password)
	}
	if config.DBName != "testdb" {
		t.Errorf("Expected DBName to be 'testdb', got '%s'", config.DBName)
	}
}

func TestDatabaseConfig_DSN(t *testing.T) {
	config := &DatabaseConfig{
		Host:     "localhost",
		Port:     "3306",
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
	}

	expected := "testuser:testpass@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := config.DSN()

	if dsn != expected {
		t.Errorf("Expected DSN to be '%s', got '%s'", expected, dsn)
	}
}

func TestNewDatabase(t *testing.T) {
	config := &DatabaseConfig{
		Host:     "localhost",
		Port:     "3306", 
		User:     "root",
		Password: "",
		DBName:   "test",
	}

	// 这个测试只验证函数不会panic，不验证实际连接
	// 因为测试环境可能没有MySQL服务器
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("NewDatabase panicked: %v", r)
		}
	}()

	_, err := NewDatabase(config)
	// 预期会有连接错误，这是正常的
	if err == nil {
		t.Log("Database connection successful (unexpected in test environment)")
	} else {
		t.Logf("Database connection failed as expected: %v", err)
	}
}