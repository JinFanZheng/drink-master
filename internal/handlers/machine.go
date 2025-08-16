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
// @Summary 获取售货机详细信息
// @Description 根据售货机ID获取售货机的详细信息
// @Tags Machine
// @Accept json
// @Produce json
// @Param id query string true "售货机ID"
// @Success 200 {object} contracts.APIResponse{data=contracts.GetMachineByIDResponse}
// @Failure 400 {object} contracts.APIResponse
// @Failure 404 {object} contracts.APIResponse
// @Router /Machine/Get [get]
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
// @Summary 检查设备是否存在
// @Description 检查指定设备ID的设备是否存在于系统中
// @Tags Machine
// @Accept json
// @Produce json
// @Param deviceId query string true "设备ID"
// @Success 200 {object} contracts.APIResponse{data=contracts.CheckDeviceExistResponse}
// @Failure 400 {object} contracts.APIResponse
// @Router /Machine/CheckDeviceExist [get]
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
// @Summary 获取售货机商品列表
// @Description 获取指定售货机的所有可售商品信息
// @Tags Machine
// @Accept json
// @Produce json
// @Param machineId query string true "售货机ID"
// @Success 200 {object} contracts.APIResponse{data=[]contracts.ProductListResponse}
// @Failure 400 {object} contracts.APIResponse
// @Router /Machine/GetProductList [get]
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
// @Summary 分页获取售货机列表
// @Description 机主权限用户分页获取其管理的售货机列表
// @Tags Machine
// @Accept json
// @Produce json
// @Param request body contracts.GetMachinePagingRequest true "分页请求"
// @Success 200 {object} contracts.APIResponse{data=contracts.GetMachinePagingResponse}
// @Failure 400 {object} contracts.APIResponse
// @Failure 401 {object} contracts.APIResponse
// @Failure 403 {object} contracts.APIResponse
// @Security BearerAuth
// @Router /Machine/GetPaging [post]
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
// @Summary 获取售货机列表
// @Description 机主权限用户获取其管理的所有售货机列表
// @Tags Machine
// @Accept json
// @Produce json
// @Success 200 {object} contracts.APIResponse{data=[]contracts.GetMachineListResponse}
// @Failure 401 {object} contracts.APIResponse
// @Failure 403 {object} contracts.APIResponse
// @Security BearerAuth
// @Router /Machine/GetList [get]
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
// @Summary 切换售货机营业状态
// @Description 机主权限用户切换指定售货机的营业状态（开启/关闭）
// @Tags Machine
// @Accept json
// @Produce json
// @Param id query string true "售货机ID"
// @Success 200 {object} contracts.APIResponse{data=contracts.OpenOrCloseBusinessResponse}
// @Failure 400 {object} contracts.APIResponse
// @Failure 401 {object} contracts.APIResponse
// @Failure 403 {object} contracts.APIResponse
// @Failure 404 {object} contracts.APIResponse
// @Security BearerAuth
// @Router /Machine/OpenOrClose [get]
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
