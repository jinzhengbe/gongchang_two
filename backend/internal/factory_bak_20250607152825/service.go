package factory

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gongChang/models"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// Register 工厂注册
func (s *Service) Register(req *RegisterRequest) error {
	// 检查用户名是否已存在
	var count int64
	s.db.Model(&Factory{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建工厂记录
	factory := &Factory{
		Name:        req.Name,
		Username:    req.Username,
		Password:    string(hashedPassword),
		Address:     req.Address,
		Contact:     req.Contact,
		Phone:       req.Phone,
		Email:       req.Email,
		License:     req.License,
		Description: req.Description,
		Status:      1, // 正常状态
	}

	return s.db.Create(factory).Error
}

// Login 工厂登录
func (s *Service) Login(req *LoginRequest) (*LoginResponse, error) {
	var factory Factory
	if err := s.db.Where("username = ?", req.Username).First(&factory).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(factory.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查工厂状态
	if factory.Status != 1 {
		return nil, errors.New("工厂账号未审核或已被禁用")
	}

	// TODO: 生成 JWT token
	token := "dummy-token" // 这里需要实现真实的 JWT token 生成

	return &LoginResponse{
		Token:   token,
		Factory: factory,
	}, nil
}

// GetRecentFactories 获取最近注册的工厂列表
func (s *Service) GetRecentFactories(limit int) ([]Factory, error) {
	var factories []Factory
	err := s.db.Order("id desc").Limit(limit).Find(&factories).Error
	if err != nil {
		return nil, err
	}
	return factories, nil
}

// GetFactories 获取工厂清单
func (s *Service) GetFactories() ([]models.Factory, error) {
	var factories []models.Factory
	if err := s.db.Find(&factories).Error; err != nil {
		return nil, err
	}
	return factories, nil
}

// GetFactoryOrders 获取工厂的订单列表
func (s *Service) GetFactoryOrders(factoryID string, req *OrderListRequest) (*OrderListResponse, error) {
	var total int64
	var orders []Order

	// 构建查询
	query := s.db.Model(&Order{}).Where("factory_id = ?", factoryID)

	// 添加筛选条件
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.StartDate != "" {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("created_at <= ?", req.EndDate)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	// 强制按id desc排序，不允许前端覆盖
	query = query.Order("id desc")

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&orders).Error; err != nil {
		return nil, err
	}

	return &OrderListResponse{
		Total:       total,
		CurrentPage: req.Page,
		PageSize:    req.PageSize,
		Orders:      orders,
	}, nil
}

// GetDesignerOrders 获取设计师的订单列表
func (s *Service) GetDesignerOrders(designerID string, req *OrderListRequest) (*OrderListResponse, error) {
	var total int64
	var orders []Order

	// 构建查询
	query := s.db.Model(&Order{}).Where("designer_id = ?", designerID)

	// 添加筛选条件
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}
	if req.StartDate != "" {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("created_at <= ?", req.EndDate)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	// 强制按id desc排序，不允许前端覆盖
	query = query.Order("id desc")

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&orders).Error; err != nil {
		return nil, err
	}

	return &OrderListResponse{
		Total:       total,
		CurrentPage: req.Page,
		PageSize:    req.PageSize,
		Orders:      orders,
	}, nil
} 