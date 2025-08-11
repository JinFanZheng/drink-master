package contracts

import "time"

// 用户相关的API契约定义

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20" example:"user123"`
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" validate:"required" example:"user123"`
	Password string `json:"password" validate:"required" example:"password123"`
}

// UserResponse 用户信息响应
type UserResponse struct {
	ID       uint      `json:"id" example:"1"`
	Username string    `json:"username" example:"user123"`
	Email    string    `json:"email" example:"user@example.com"`
	CreateAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdateAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt time.Time    `json:"expires_at" example:"2023-01-02T00:00:00Z"`
	User      UserResponse `json:"user"`
}

// TokenClaims JWT Token声明
type TokenClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}