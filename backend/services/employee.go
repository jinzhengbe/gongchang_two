package services

import (
	"errors"
	"time"
	"gorm.io/gorm"
	"gongChang/models"
)

type EmployeeService struct {
	db *gorm.DB
}

func NewEmployeeService(db *gorm.DB) *EmployeeService {
	return &EmployeeService{db: db}
}

// CreateEmployee 创建职工
func (s *EmployeeService) CreateEmployee(factoryID string, req *models.CreateEmployeeRequest) (*models.FactoryEmployee, error) {
	// 验证工厂是否存在
	var factory models.FactoryProfile
	if err := s.db.Where("user_id = ?", factoryID).First(&factory).Error; err != nil {
		return nil, errors.New("工厂不存在")
	}

	employee := &models.FactoryEmployee{
		Name:       req.Name,
		Position:   req.Position,
		Grade:      req.Grade,
		WorkYears:  req.WorkYears,
		FactoryID:  factoryID,
		HireDate:   req.HireDate,
		Phone:      req.Phone,
		Email:      req.Email,
		Department: req.Department,
		Salary:     req.Salary,
		Status:     req.Status,
	}

	if employee.Status == "" {
		employee.Status = models.EmployeeStatusActive
	}

	if err := s.db.Create(employee).Error; err != nil {
		return nil, err
	}

	return employee, nil
}

// GetEmployeesByFactory 获取工厂职工列表
func (s *EmployeeService) GetEmployeesByFactory(factoryID string, page, pageSize int, status string, department string) (*models.EmployeeListResponse, error) {
	var employees []models.FactoryEmployee
	var total int64

	query := s.db.Where("factory_id = ?", factoryID).Preload("Factory")

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 部门筛选
	if department != "" {
		query = query.Where("department = ?", department)
	}

	// 获取总数
	if err := query.Model(&models.FactoryEmployee{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&employees).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	employeeResponses := make([]models.EmployeeResponse, len(employees))
	for i, emp := range employees {
		employeeResponses[i] = models.EmployeeResponse{
			ID:         emp.ID,
			Name:       emp.Name,
			Position:   emp.Position,
			Grade:      emp.Grade,
			WorkYears:  emp.WorkYears,
			FactoryID:  emp.FactoryID,
			HireDate:   emp.HireDate,
			Phone:      emp.Phone,
			Email:      emp.Email,
			Department: emp.Department,
			Salary:     emp.Salary,
			Status:     emp.Status,
			CreatedAt:  emp.CreatedAt,
			UpdatedAt:  emp.UpdatedAt,
			Factory:    &emp.Factory,
		}
	}

	return &models.EmployeeListResponse{
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		Employees: employeeResponses,
	}, nil
}

// GetEmployeeByID 根据ID获取职工
func (s *EmployeeService) GetEmployeeByID(factoryID string, employeeID uint) (*models.FactoryEmployee, error) {
	var employee models.FactoryEmployee
	err := s.db.Where("id = ? AND factory_id = ?", employeeID, factoryID).Preload("Factory").First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

// UpdateEmployee 更新职工信息
func (s *EmployeeService) UpdateEmployee(factoryID string, employeeID uint, req *models.UpdateEmployeeRequest) (*models.FactoryEmployee, error) {
	employee, err := s.GetEmployeeByID(factoryID, employeeID)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Position != nil {
		updates["position"] = *req.Position
	}
	if req.Grade != nil {
		updates["grade"] = *req.Grade
	}
	if req.WorkYears != nil {
		updates["work_years"] = *req.WorkYears
	}
	if req.HireDate != nil {
		updates["hire_date"] = *req.HireDate
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Department != nil {
		updates["department"] = *req.Department
	}
	if req.Salary != nil {
		updates["salary"] = *req.Salary
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	updates["updated_at"] = time.Now()

	if err := s.db.Model(employee).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 重新查询获取更新后的数据
	return s.GetEmployeeByID(factoryID, employeeID)
}

// DeleteEmployee 删除职工
func (s *EmployeeService) DeleteEmployee(factoryID string, employeeID uint) error {
	employee, err := s.GetEmployeeByID(factoryID, employeeID)
	if err != nil {
		return err
	}

	return s.db.Delete(employee).Error
}

// GetEmployeeStatistics 获取职工统计
func (s *EmployeeService) GetEmployeeStatistics(factoryID string) (*models.EmployeeStatistics, error) {
	var stats models.EmployeeStatistics

	// 总职工数
	if err := s.db.Model(&models.FactoryEmployee{}).Where("factory_id = ?", factoryID).Count(&stats.TotalEmployees).Error; err != nil {
		return nil, err
	}

	// 在职职工数
	if err := s.db.Model(&models.FactoryEmployee{}).Where("factory_id = ? AND status = ?", factoryID, models.EmployeeStatusActive).Count(&stats.ActiveEmployees).Error; err != nil {
		return nil, err
	}

	// 离职职工数
	if err := s.db.Model(&models.FactoryEmployee{}).Where("factory_id = ? AND status = ?", factoryID, models.EmployeeStatusInactive).Count(&stats.InactiveEmployees).Error; err != nil {
		return nil, err
	}

	// 平均工龄
	var avgWorkYears *float64
	if err := s.db.Model(&models.FactoryEmployee{}).Where("factory_id = ?", factoryID).Select("AVG(work_years)").Scan(&avgWorkYears).Error; err != nil {
		return nil, err
	}
	
	if avgWorkYears != nil {
		stats.AverageWorkYears = *avgWorkYears
	} else {
		stats.AverageWorkYears = 0.0
	}

	// 部门统计
	var departmentStats []struct {
		Department string `json:"department"`
		Count      int64  `json:"count"`
	}
	if err := s.db.Model(&models.FactoryEmployee{}).
		Where("factory_id = ? AND department IS NOT NULL", factoryID).
		Select("department, COUNT(*) as count").
		Group("department").
		Scan(&departmentStats).Error; err != nil {
		return nil, err
	}

	stats.DepartmentStats = make(map[string]int64)
	for _, dept := range departmentStats {
		stats.DepartmentStats[dept.Department] = dept.Count
	}

	return &stats, nil
}

// SearchEmployees 搜索职工
func (s *EmployeeService) SearchEmployees(factoryID string, keyword string, page, pageSize int) (*models.EmployeeListResponse, error) {
	var employees []models.FactoryEmployee
	var total int64

	query := s.db.Where("factory_id = ?", factoryID).Preload("Factory")

	if keyword != "" {
		query = query.Where("name LIKE ? OR position LIKE ? OR department LIKE ?", 
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err := query.Model(&models.FactoryEmployee{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&employees).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	employeeResponses := make([]models.EmployeeResponse, len(employees))
	for i, emp := range employees {
		employeeResponses[i] = models.EmployeeResponse{
			ID:         emp.ID,
			Name:       emp.Name,
			Position:   emp.Position,
			Grade:      emp.Grade,
			WorkYears:  emp.WorkYears,
			FactoryID:  emp.FactoryID,
			HireDate:   emp.HireDate,
			Phone:      emp.Phone,
			Email:      emp.Email,
			Department: emp.Department,
			Salary:     emp.Salary,
			Status:     emp.Status,
			CreatedAt:  emp.CreatedAt,
			UpdatedAt:  emp.UpdatedAt,
			Factory:    &emp.Factory,
		}
	}

	return &models.EmployeeListResponse{
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		Employees: employeeResponses,
	}, nil
} 