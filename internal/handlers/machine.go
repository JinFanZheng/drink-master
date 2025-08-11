package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ddteam/drink-master/internal/contracts"
	"github.com/ddteam/drink-master/internal/services"
)

// MachineHandler 售货机处理器 (对应MobileAPI MachineController)
type MachineHandler struct {
	*BaseHandler
	machineService services.MachineServiceInterface
}

// NewMachineHandler 创建售货机处理器
func NewMachineHandler(db *gorm.DB) *MachineHandler {
	return &MachineHandler{
		BaseHandler:    NewBaseHandler(db),
		machineService: services.NewMachineService(db),
	}
}

// Get 获取售货机信息
// GET /api/Machine/Get?id=machine_id
func (h *MachineHandler) Get(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		h.ErrorResponse(c, http.StatusBadRequest, contracts.ErrorCodeValidation, "machine id required")
		return
	}

	machine, err := h.machineService.GetMachineByID(id)
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	if machine == nil {
		h.NotFoundResponse(c, "machine not found")
		return
	}

	h.SuccessResponse(c, machine)
}

// CheckDeviceExist 检查设备是否存在
// GET /api/Machine/CheckDeviceExist?deviceId=device_id
func (h *MachineHandler) CheckDeviceExist(c *gin.Context) {
	var req contracts.CheckDeviceExistRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	exists, err := h.machineService.CheckDeviceExist(req.DeviceID)
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	response := contracts.CheckDeviceExistResponse{
		Exists: exists,
	}

	h.SuccessResponse(c, response)
}

// GetProductList 获取售货机商品列表
// GET /api/Machine/GetProductList?machineId=machine_id
func (h *MachineHandler) GetProductList(c *gin.Context) {
	machineID := c.Query("machineId")
	if machineID == "" {
		h.ErrorResponse(c, http.StatusBadRequest, contracts.ErrorCodeValidation, "machineId required")
		return
	}

	products, err := h.machineService.GetProductList(machineID)
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	// 返回VendingMachine格式
	if len(products) == 0 {
		h.SuccessResponse(c, []interface{}{})
		return
	}

	h.SuccessResponse(c, products)
}

// GetPaging 分页获取售货机列表（需要机主权限）
// POST /api/Machine/GetPaging
func (h *MachineHandler) GetPaging(c *gin.Context) {
	// 检查机主权限
	if !h.IsMachineOwner(c) {
		h.ForbiddenResponse(c, "machine owner permission required")
		return
	}

	// 获取机主ID
	machineOwnerID, exists := h.GetMachineOwnerID(c)
	if !exists {
		h.UnauthorizedResponse(c, "machine owner id not found")
		return
	}

	var req contracts.GetMachinePagingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	// 设置机主ID
	req.MachineOwnerID = machineOwnerID

	result, err := h.machineService.GetMachinePaging(req)
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	h.SuccessResponse(c, result)
}

// GetList 获取售货机列表（需要机主权限）
// GET /api/Machine/GetList
func (h *MachineHandler) GetList(c *gin.Context) {
	// 检查机主权限
	if !h.IsMachineOwner(c) {
		h.ForbiddenResponse(c, "machine owner permission required")
		return
	}

	// 获取机主ID
	machineOwnerID, exists := h.GetMachineOwnerID(c)
	if !exists {
		h.UnauthorizedResponse(c, "machine owner id not found")
		return
	}

	machines, err := h.machineService.GetMachineList(machineOwnerID)
	if err != nil {
		h.InternalErrorResponse(c, err)
		return
	}

	h.SuccessResponse(c, machines)
}

// OpenOrCloseBusiness 开启或关闭营业状态（需要机主权限）
// GET /api/Machine/OpenOrClose?id=machine_id
func (h *MachineHandler) OpenOrCloseBusiness(c *gin.Context) {
	// 检查机主权限
	if !h.IsMachineOwner(c) {
		h.ForbiddenResponse(c, "machine owner permission required")
		return
	}

	// 获取机主ID
	machineOwnerID, exists := h.GetMachineOwnerID(c)
	if !exists {
		h.UnauthorizedResponse(c, "machine owner id not found")
		return
	}

	var req contracts.OpenOrCloseBusinessRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.machineService.OpenOrCloseBusiness(req.ID, machineOwnerID)
	if err != nil {
		if err.Error() == "machine not found" {
			h.NotFoundResponse(c, "machine not found")
			return
		}
		if err.Error() == "permission denied: not machine owner" {
			h.ForbiddenResponse(c, "permission denied")
			return
		}
		h.InternalErrorResponse(c, err)
		return
	}

	h.SuccessResponseWithMessage(c, result, result.Message)
}
