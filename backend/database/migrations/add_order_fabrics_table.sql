-- 创建订单-布料关联表
CREATE TABLE IF NOT EXISTS order_fabrics (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL,
    fabric_id BIGINT UNSIGNED NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- 添加外键约束
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (fabric_id) REFERENCES fabrics(id) ON DELETE CASCADE,
    
    -- 添加唯一约束，防止重复关联
    UNIQUE KEY unique_order_fabric (order_id, fabric_id),
    
    -- 添加索引
    INDEX idx_order_id (order_id),
    INDEX idx_fabric_id (fabric_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci; 