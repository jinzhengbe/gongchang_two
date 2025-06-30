package controllers

import (
	"net/http"
	"strconv"
	"gongChang/models"
	"gongChang/services"
	"github.com/gin-gonic/gin"
)

type EmployeeController struct {
	employeeService *services.EmployeeService
}

func NewEmployeeController(employeeService *services.EmployeeService) *EmployeeController {
	return &EmployeeController{
		employeeService: employeeService,
	}
}

// CreateEmployee 创建职工
func (c *EmployeeController) CreateEmployee(ctx *gin.Context) {
	factoryID := ctx.GetString("user_id")
	if factoryID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var req models.CreateEmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	employee, err := c.employeeService.CreateEmployee(factoryID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "职工创建成功",
		"employee": employee,
	})
}

// GetEmployees 获取工厂职工列表
func (c *EmployeeController) GetEmployees(ctx *gin.Context) {
	factoryID := ctx.GetString("user_id")
	if factoryID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 获取查询参数
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	status := ctx.Query("status")
	department := ctx.Query("department")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// 限制最大页面大小
	if pageSize > 100 {
		pageSize = 100
	}

	result, err := c.employeeService.GetEmployeesByFactory(factoryID, page, pageSize, status, department)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetEmployee 获取单个职工信息
func (c *EmployeeController) GetEmployee(ctx *gin.Context) {
	factoryID := ctx.GetString("user_id")
	if factoryID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	employeeIDStr := ctx.Param("id")
	employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的职工ID"})
		return
	}

	employee, err := c.employeeService.GetEmployeeByID(factoryID, uint(employeeID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "职工不存在"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"employee": employee,
	})
}

// UpdateEmployee 更新职工信息
func (c *EmployeeController) UpdateEmployee(ctx *gin.Context) {
	factoryID := ctx.GetString("user_id")
	if factoryID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	employeeIDStr := ctx.Param("id")
	employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的职工ID"})
		return
	}

	var req models.UpdateEmployeeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	employee, err := c.employeeService.UpdateEmployee(factoryID, uint(employeeID), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "职工信息更新成功",
		"employee": employee,
	})
}

// DeleteEmployee 删除职工
func (c *EmployeeController) DeleteEmployee(ctx *gin.Context) {
	factoryID := ctx.GetString("user_id")
	if factoryID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	employeeIDStr := ctx.Param("id")
	employeeID, err := strconv.ParseUint(employeeIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的职工ID"})
		return
	}

	err = c.employeeService.DeleteEmployee(factoryID, uint(employeeID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "职工删除成功",
	})
}

// GetEmployeeStatistics 获取职工统计
func (c *EmployeeController) GetEmployeeStatistics(ctx *gin.Context) {
	factoryID := ctx.GetString("user_id")
	if factoryID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	stats, err := c.employeeService.GetEmployeeStatistics(factoryID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"statistics": stats,
	})
}

// SearchEmployees 搜索职工
func (c *EmployeeController) SearchEmployees(ctx *gin.Context) {
	factoryID := ctx.GetString("user_id")
	if factoryID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	keyword := ctx.Query("q")
	if keyword == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}

	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	if pageSize > 100 {
		pageSize = 100
	}

	result, err := c.employeeService.SearchEmployees(factoryID, keyword, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
} 