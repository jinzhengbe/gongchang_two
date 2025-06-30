package models

import (
	"time"
	"gorm.io/gorm"
)

// EmployeeStatus 职工状态
type EmployeeStatus string

const (
	EmployeeStatusActive   EmployeeStatus = "active"   // 在职
	EmployeeStatusInactive EmployeeStatus = "inactive" // 离职
)

// FactoryEmployee 工厂职工模型
type FactoryEmployee struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Name       string         `json:"name" gorm:"type:varchar(100);not null;comment:职工姓名"`
	Position   string         `json:"position" gorm:"type:varchar(100);not null;comment:职位"`
	Grade      *string        `json:"grade" gorm:"type:varchar(50);comment:年级/级别"`
	WorkYears  int            `json:"work_years" gorm:"default:0;comment:工龄(年)"`
	FactoryID  string         `json:"factory_id" gorm:"type:varchar(191);not null;index;comment:工厂ID"`
	HireDate   time.Time      `json:"hire_date" gorm:"type:date;not null;comment:入职时间"`
	Phone      *string        `json:"phone" gorm:"type:varchar(20);comment:联系电话"`
	Email      *string        `json:"email" gorm:"type:varchar(100);comment:邮箱"`
	Department *string        `json:"department" gorm:"type:varchar(100);comment:部门"`
	Salary     *float64       `json:"salary" gorm:"type:decimal(10,2);comment:薪资"`
	Status     EmployeeStatus `json:"status" gorm:"type:varchar(20);default:'active';comment:状态"`
	CreatedAt  *time.Time     `json:"created_at" gorm:"autoCreateTime:false"`
	UpdatedAt  *time.Time     `json:"updated_at" gorm:"autoUpdateTime:false"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	
	// 关联关系
	Factory FactoryProfile `json:"factory" gorm:"foreignKey:FactoryID;references:UserID"`
}

// TableName 指定表名
func (FactoryEmployee) TableName() string {
	return "factory_employees"
}

// CreateEmployeeRequest 创建职工请求
type CreateEmployeeRequest struct {
	Name       string         `json:"name" binding:"required"`
	Position   string         `json:"position" binding:"required"`
	Grade      *string        `json:"grade"`
	WorkYears  int            `json:"work_years"`
	HireDate   time.Time      `json:"hire_date" binding:"required"`
	Phone      *string        `json:"phone"`
	Email      *string        `json:"email"`
	Department *string        `json:"department"`
	Salary     *float64       `json:"salary"`
	Status     EmployeeStatus `json:"status"`
}

// UpdateEmployeeRequest 更新职工请求
type UpdateEmployeeRequest struct {
	Name       *string         `json:"name"`
	Position   *string         `json:"position"`
	Grade      *string         `json:"grade"`
	WorkYears  *int            `json:"work_years"`
	HireDate   *time.Time      `json:"hire_date"`
	Phone      *string         `json:"phone"`
	Email      *string         `json:"email"`
	Department *string         `json:"department"`
	Salary     *float64        `json:"salary"`
	Status     *EmployeeStatus `json:"status"`
}

// EmployeeResponse 职工响应
type EmployeeResponse struct {
	ID         uint           `json:"id"`
	Name       string         `json:"name"`
	Position   string         `json:"position"`
	Grade      *string        `json:"grade"`
	WorkYears  int            `json:"work_years"`
	FactoryID  string         `json:"factory_id"`
	HireDate   time.Time      `json:"hire_date"`
	Phone      *string        `json:"phone"`
	Email      *string        `json:"email"`
	Department *string        `json:"department"`
	Salary     *float64       `json:"salary"`
	Status     EmployeeStatus `json:"status"`
	CreatedAt  *time.Time     `json:"created_at"`
	UpdatedAt  *time.Time     `json:"updated_at"`
	
	// 关联数据
	Factory *FactoryProfile `json:"factory,omitempty"`
}

// EmployeeListResponse 职工列表响应
type EmployeeListResponse struct {
	Total     int64            `json:"total"`
	Page      int              `json:"page"`
	PageSize  int              `json:"page_size"`
	Employees []EmployeeResponse `json:"employees"`
}

// EmployeeStatistics 职工统计
type EmployeeStatistics struct {
	TotalEmployees    int64 `json:"total_employees"`
	ActiveEmployees   int64 `json:"active_employees"`
	InactiveEmployees int64 `json:"inactive_employees"`
	AverageWorkYears  float64 `json:"average_work_years"`
	DepartmentStats   map[string]int64 `json:"department_stats"`
} 