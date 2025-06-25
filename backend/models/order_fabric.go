package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// OrderFabric 订单-布料关联模型
type OrderFabric struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	OrderID   uint      `json:"order_id" gorm:"index"`
	FabricID  uint      `json:"fabric_id" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// 关联关系
	Order  Order  `json:"order" gorm:"foreignKey:OrderID"`
	Fabric Fabric `json:"fabric" gorm:"foreignKey:FabricID"`
}

// AddFabricToOrderRequest 添加布料到订单的请求
type AddFabricToOrderRequest struct {
	OrderID uint `json:"order_id" binding:"required"`
	// 布料信息
	Name         string  `json:"name" binding:"required"`
	Category     string  `json:"category"`
	Material     string  `json:"material"`
	Color        string  `json:"color"`
	Pattern      string  `json:"pattern"`
	Weight       float64 `json:"weight"`
	Width        float64 `json:"width"`
	Price        float64 `json:"price"`
	Unit         string  `json:"unit"`
	Stock        int     `json:"stock"`
	MinOrder     int     `json:"min_order"`
	Description  string  `json:"description"`
	ImageURL     string  `json:"image_url"`
	ThumbnailURL string  `json:"thumbnail_url"`
	Tags         string  `json:"tags"`
}

// AddFabricToOrderResponse 添加布料到订单的响应
type AddFabricToOrderResponse struct {
	Message           string  `json:"message"`
	Fabric            *Fabric `json:"fabric"`
	OrderID           uint    `json:"order_id"`
	AssociationCreated bool   `json:"association_created"`
}

// FabricIDList 布料ID列表工具函数
type FabricIDList []uint

// AddFabricID 添加布料ID到列表
func (f *FabricIDList) AddFabricID(fabricID uint) {
	*f = append(*f, fabricID)
}

// RemoveFabricID 从列表中移除布料ID
func (f *FabricIDList) RemoveFabricID(fabricID uint) {
	for i, id := range *f {
		if id == fabricID {
			*f = append((*f)[:i], (*f)[i+1:]...)
			break
		}
	}
}

// ContainsFabricID 检查是否包含指定布料ID
func (f *FabricIDList) ContainsFabricID(fabricID uint) bool {
	for _, id := range *f {
		if id == fabricID {
			return true
		}
	}
	return false
}

// ToJSONString 转换为JSON字符串
func (f *FabricIDList) ToJSONString() string {
	if len(*f) == 0 {
		return "[]"
	}
	jsonBytes, _ := json.Marshal(*f)
	return string(jsonBytes)
}

// FromJSONString 从JSON字符串解析
func (f *FabricIDList) FromJSONString(jsonStr string) error {
	if jsonStr == "" || jsonStr == "[]" {
		*f = make(FabricIDList, 0)
		return nil
	}
	return json.Unmarshal([]byte(jsonStr), f)
}

// FromCommaString 从逗号分隔的字符串解析
func (f *FabricIDList) FromCommaString(commaStr string) error {
	if commaStr == "" {
		*f = make(FabricIDList, 0)
		return nil
	}
	
	parts := strings.Split(commaStr, ",")
	*f = make(FabricIDList, 0, len(parts))
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			if id, err := strconv.ParseUint(part, 10, 32); err == nil {
				*f = append(*f, uint(id))
			}
		}
	}
	return nil
}

// ToCommaString 转换为逗号分隔的字符串
func (f *FabricIDList) ToCommaString() string {
	if len(*f) == 0 {
		return ""
	}
	
	strParts := make([]string, len(*f))
	for i, id := range *f {
		strParts[i] = strconv.FormatUint(uint64(id), 10)
	}
	return strings.Join(strParts, ",")
} 