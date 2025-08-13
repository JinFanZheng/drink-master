# 错误处理模式

## 错误类型定义

### 自定义错误类型
```go
// internal/errors/errors.go
type AppError struct {
    Code       string
    Message    string
    StatusCode int
    Cause      error
}

func (e *AppError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Cause)
    }
    return e.Message
}

func (e *AppError) Unwrap() error {
    return e.Cause
}

// 预定义错误
var (
    ErrNotFound = &AppError{
        Code:       "NOT_FOUND",
        Message:    "Resource not found",
        StatusCode: 404,
    }
    
    ErrUnauthorized = &AppError{
        Code:       "UNAUTHORIZED",
        Message:    "Unauthorized access",
        StatusCode: 401,
    }
    
    ErrInvalidInput = &AppError{
        Code:       "INVALID_INPUT",
        Message:    "Invalid input parameters",
        StatusCode: 400,
    }
)
```

## Service层错误处理

### 包装错误
```go
func (s *MachineService) GetMachine(id uint) (*models.Machine, error) {
    var machine models.Machine
    err := s.db.First(&machine, id).Error
    
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, &AppError{
                Code:       "MACHINE_NOT_FOUND",
                Message:    fmt.Sprintf("Machine with ID %d not found", id),
                StatusCode: 404,
                Cause:      err,
            }
        }
        // 包装未知错误
        return nil, &AppError{
            Code:       "DATABASE_ERROR",
            Message:    "Failed to fetch machine",
            StatusCode: 500,
            Cause:      err,
        }
    }
    
    return &machine, nil
}
```

### 错误链处理
```go
func (s *OrderService) CreateOrder(req *CreateOrderRequest) (*models.Order, error) {
    // 验证库存
    if err := s.checkInventory(req.ProductID, req.Quantity); err != nil {
        return nil, fmt.Errorf("inventory check failed: %w", err)
    }
    
    // 验证用户
    if err := s.validateUser(req.UserID); err != nil {
        return nil, fmt.Errorf("user validation failed: %w", err)
    }
    
    // 创建订单
    order, err := s.createOrderRecord(req)
    if err != nil {
        return nil, fmt.Errorf("create order failed: %w", err)
    }
    
    return order, nil
}
```

## Handler层错误响应

### 统一错误响应
```go
func HandleError(c *gin.Context, err error) {
    // 检查是否是自定义错误
    var appErr *AppError
    if errors.As(err, &appErr) {
        c.JSON(appErr.StatusCode, gin.H{
            "success": false,
            "error": gin.H{
                "code":    appErr.Code,
                "message": appErr.Message,
            },
        })
        
        // 记录错误日志
        if appErr.StatusCode >= 500 {
            log.Printf("Internal error: %+v", appErr)
        }
        return
    }
    
    // 检查是否是验证错误
    var validationErr validator.ValidationErrors
    if errors.As(err, &validationErr) {
        c.JSON(400, gin.H{
            "success": false,
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": "Validation failed",
                "details": formatValidationErrors(validationErr),
            },
        })
        return
    }
    
    // 默认内部错误
    log.Printf("Unexpected error: %+v", err)
    c.JSON(500, gin.H{
        "success": false,
        "error": gin.H{
            "code":    "INTERNAL_ERROR",
            "message": "An internal error occurred",
        },
    })
}
```

### Handler使用示例
```go
func GetMachine(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        HandleError(c, &AppError{
            Code:       "INVALID_ID",
            Message:    "Invalid machine ID",
            StatusCode: 400,
            Cause:      err,
        })
        return
    }
    
    machine, err := machineService.GetMachine(uint(id))
    if err != nil {
        HandleError(c, err)
        return
    }
    
    c.JSON(200, gin.H{
        "success": true,
        "data":    machine,
    })
}
```

## 错误恢复

### Panic恢复中间件
```go
func RecoveryMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // 记录panic详情
                log.Printf("Panic recovered: %v\n%s", err, debug.Stack())
                
                c.JSON(500, gin.H{
                    "success": false,
                    "error": gin.H{
                        "code":    "PANIC",
                        "message": "An unexpected error occurred",
                    },
                })
                c.Abort()
            }
        }()
        c.Next()
    }
}
```

### 重试机制
```go
func RetryOperation(operation func() error, maxAttempts int) error {
    var lastErr error
    
    for i := 0; i < maxAttempts; i++ {
        if err := operation(); err != nil {
            lastErr = err
            
            // 不重试的错误类型
            var appErr *AppError
            if errors.As(err, &appErr) {
                if appErr.StatusCode < 500 {
                    return err // 客户端错误不重试
                }
            }
            
            // 指数退避
            waitTime := time.Duration(math.Pow(2, float64(i))) * time.Second
            time.Sleep(waitTime)
            continue
        }
        
        return nil // 成功
    }
    
    return fmt.Errorf("operation failed after %d attempts: %w", maxAttempts, lastErr)
}

// 使用示例
err := RetryOperation(func() error {
    return callExternalAPI()
}, 3)
```

## 错误日志

### 结构化日志
```go
type ErrorLogger struct {
    logger *log.Logger
}

func (e *ErrorLogger) LogError(ctx context.Context, err error, extra ...interface{}) {
    requestID := ctx.Value("request_id")
    userID := ctx.Value("user_id")
    
    logEntry := map[string]interface{}{
        "timestamp":  time.Now().Unix(),
        "request_id": requestID,
        "user_id":    userID,
        "error":      err.Error(),
    }
    
    // 添加额外信息
    for i := 0; i < len(extra); i += 2 {
        if i+1 < len(extra) {
            key := fmt.Sprint(extra[i])
            value := extra[i+1]
            logEntry[key] = value
        }
    }
    
    // 判断错误级别
    var appErr *AppError
    if errors.As(err, &appErr) {
        if appErr.StatusCode >= 500 {
            e.logger.Printf("ERROR: %+v", logEntry)
        } else {
            e.logger.Printf("WARN: %+v", logEntry)
        }
    } else {
        e.logger.Printf("ERROR: %+v", logEntry)
    }
}
```

## 错误监控

### 错误率统计
```go
var (
    errorCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "api_errors_total",
            Help: "Total number of API errors",
        },
        []string{"code", "endpoint"},
    )
)

func init() {
    prometheus.MustRegister(errorCounter)
}

func RecordError(code string, endpoint string) {
    errorCounter.WithLabelValues(code, endpoint).Inc()
}
```

### 告警规则
```yaml
# Prometheus告警规则
groups:
  - name: api_errors
    rules:
      - alert: HighErrorRate
        expr: rate(api_errors_total[5m]) > 0.05
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"
          description: "Error rate is above 5% for 5 minutes"
```

## 最佳实践

### 1. 错误要有上下文
```go
// ❌ 不好
return errors.New("failed")

// ✅ 好
return fmt.Errorf("failed to create order for user %d: %w", userID, err)
```

### 2. 区分预期和意外错误
```go
// 预期错误（业务错误）
if user.Balance < orderAmount {
    return &AppError{
        Code:       "INSUFFICIENT_BALANCE",
        Message:    "Insufficient balance",
        StatusCode: 400,
    }
}

// 意外错误（系统错误）
if err := db.Save(&order).Error; err != nil {
    return fmt.Errorf("database error: %w", err)
}
```

### 3. 不要忽略错误
```go
// ❌ 不好
_ = doSomething()

// ✅ 好
if err := doSomething(); err != nil {
    // 至少记录日志
    log.Printf("doSomething failed: %v", err)
}
```

### 4. 提供有用的错误信息
```go
// ❌ 不好
return &AppError{
    Message: "Error",
}

// ✅ 好
return &AppError{
    Code:    "PRODUCT_OUT_OF_STOCK",
    Message: fmt.Sprintf("Product %s is out of stock, available: %d, requested: %d", 
        product.Name, product.Stock, requested),
}
```