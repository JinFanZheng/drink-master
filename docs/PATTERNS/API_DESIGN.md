# API设计模式

## RESTful API标准

### 路由命名规范
```go
// ✅ 正确
router.GET("/api/machines", handler.GetMachines)       // 复数，获取列表
router.GET("/api/machines/:id", handler.GetMachine)    // 获取单个
router.POST("/api/machines", handler.CreateMachine)    // 创建
router.PUT("/api/machines/:id", handler.UpdateMachine) // 更新
router.DELETE("/api/machines/:id", handler.DeleteMachine) // 删除

// ❌ 错误
router.GET("/api/getMachine", handler.GetMachine)      // 动词在URL中
router.POST("/api/machine/create", handler.Create)     // 动词冗余
```

### 请求/响应格式

#### 标准响应结构
```go
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorInfo  `json:"error,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorInfo struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

type Meta struct {
    Page       int `json:"page,omitempty"`
    PageSize   int `json:"page_size,omitempty"`
    Total      int `json:"total,omitempty"`
    TotalPages int `json:"total_pages,omitempty"`
}
```

#### 分页请求
```go
// 请求参数
type PaginationRequest struct {
    Page     int `form:"page" binding:"min=1" default:"1"`
    PageSize int `form:"page_size" binding:"min=1,max=100" default:"20"`
    Sort     string `form:"sort" default:"-created_at"`
}

// 使用示例
func GetMachines(c *gin.Context) {
    var req PaginationRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(400, ErrorResponse(err))
        return
    }
    // 处理逻辑...
}
```

## 错误处理

### 错误码设计
```go
const (
    // 业务错误码
    ErrCodeInvalidInput   = "INVALID_INPUT"
    ErrCodeNotFound       = "NOT_FOUND"
    ErrCodeUnauthorized   = "UNAUTHORIZED"
    ErrCodeForbidden      = "FORBIDDEN"
    ErrCodeConflict       = "CONFLICT"
    ErrCodeInternal       = "INTERNAL_ERROR"
)

// HTTP状态码映射
var errorStatusMap = map[string]int{
    ErrCodeInvalidInput: 400,
    ErrCodeNotFound:     404,
    ErrCodeUnauthorized: 401,
    ErrCodeForbidden:    403,
    ErrCodeConflict:     409,
    ErrCodeInternal:     500,
}
```

### 统一错误处理
```go
func HandleError(c *gin.Context, err error) {
    var apiErr *APIError
    if errors.As(err, &apiErr) {
        c.JSON(apiErr.StatusCode, Response{
            Success: false,
            Error: &ErrorInfo{
                Code:    apiErr.Code,
                Message: apiErr.Message,
            },
        })
        return
    }
    
    // 默认内部错误
    c.JSON(500, Response{
        Success: false,
        Error: &ErrorInfo{
            Code:    ErrCodeInternal,
            Message: "Internal server error",
        },
    })
}
```

## 参数验证

### 使用Gin的binding
```go
type CreateMachineRequest struct {
    Name         string  `json:"name" binding:"required,min=1,max=100"`
    Location     string  `json:"location" binding:"required"`
    Capacity     int     `json:"capacity" binding:"required,min=1,max=1000"`
    Status       string  `json:"status" binding:"oneof=active inactive maintenance"`
}

func CreateMachine(c *gin.Context) {
    var req CreateMachineRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, Response{
            Success: false,
            Error: &ErrorInfo{
                Code:    ErrCodeInvalidInput,
                Message: err.Error(),
            },
        })
        return
    }
    // 业务逻辑...
}
```

### 自定义验证
```go
// 注册自定义验证器
func init() {
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterValidation("phone", validatePhone)
    }
}

func validatePhone(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
    return matched
}

// 使用
type UserRequest struct {
    Phone string `json:"phone" binding:"required,phone"`
}
```

## 认证与授权

### JWT中间件
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            token = c.Query("token")
        }
        
        if token == "" {
            c.JSON(401, Response{
                Success: false,
                Error: &ErrorInfo{
                    Code:    ErrCodeUnauthorized,
                    Message: "Missing token",
                },
            })
            c.Abort()
            return
        }
        
        claims, err := ValidateToken(token)
        if err != nil {
            c.JSON(401, Response{
                Success: false,
                Error: &ErrorInfo{
                    Code:    ErrCodeUnauthorized,
                    Message: "Invalid token",
                },
            })
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Set("user_role", claims.Role)
        c.Next()
    }
}
```

## API版本管理

### URL版本
```go
v1 := router.Group("/api/v1")
{
    v1.GET("/machines", v1Handler.GetMachines)
}

v2 := router.Group("/api/v2")
{
    v2.GET("/machines", v2Handler.GetMachines)
}
```

### Header版本
```go
func VersionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        version := c.GetHeader("API-Version")
        if version == "" {
            version = "v1" // 默认版本
        }
        c.Set("api_version", version)
        c.Next()
    }
}
```

## 最佳实践

### 1. 幂等性设计
```go
// PUT和DELETE应该是幂等的
func UpdateMachine(c *gin.Context) {
    id := c.Param("id")
    
    // 使用乐观锁防止并发更新
    var req UpdateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        HandleError(c, err)
        return
    }
    
    err := db.Model(&Machine{}).
        Where("id = ? AND version = ?", id, req.Version).
        Updates(req).Error
    
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            HandleError(c, ErrVersionConflict)
            return
        }
        HandleError(c, err)
        return
    }
}
```

### 2. 限流保护
```go
// 使用中间件限流
rateLimiter := rate.NewLimiter(100, 10) // 100 req/s, burst 10

func RateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !rateLimiter.Allow() {
            c.JSON(429, Response{
                Success: false,
                Error: &ErrorInfo{
                    Code:    "RATE_LIMIT",
                    Message: "Too many requests",
                },
            })
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 3. 请求追踪
```go
func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        c.Set("request_id", requestID)
        c.Header("X-Request-ID", requestID)
        c.Next()
    }
}
```