-- 订单搜索相关索引
-- 执行此脚本以优化订单搜索性能

-- 1. 订单搜索复合索引
CREATE INDEX IF NOT EXISTS idx_orders_search ON orders (
    title, 
    status, 
    created_at
);

-- 2. 订单面料关联索引
CREATE INDEX IF NOT EXISTS idx_orders_fabric ON orders (
    fabric,
    fabrics
);

-- 3. 订单时间范围索引
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders (
    created_at DESC
);

-- 4. 订单状态索引
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders (
    status
);

-- 5. 订单用户关联索引
CREATE INDEX IF NOT EXISTS idx_orders_factory_id ON orders (
    factory_id
);

CREATE INDEX IF NOT EXISTS idx_orders_designer_id ON orders (
    designer_id
);

-- 6. 工厂名称搜索索引
CREATE INDEX IF NOT EXISTS idx_factory_profiles_company_name ON factory_profiles (
    company_name
);

-- 7. 订单描述全文搜索索引（如果数据库支持）
-- MySQL 5.7+ 支持全文搜索
-- ALTER TABLE orders ADD FULLTEXT INDEX idx_orders_fulltext (title, description, fabric);

-- 8. 订单号索引（如果使用订单号）
-- CREATE INDEX IF NOT EXISTS idx_orders_order_no ON orders (order_no);

-- 9. 复合查询索引
CREATE INDEX IF NOT EXISTS idx_orders_composite_search ON orders (
    factory_id,
    status,
    created_at DESC
);

CREATE INDEX IF NOT EXISTS idx_orders_composite_designer ON orders (
    designer_id,
    status,
    created_at DESC
);

-- 10. 删除标记索引（软删除）
CREATE INDEX IF NOT EXISTS idx_orders_deleted_at ON orders (
    deleted_at
);

-- 查看索引创建结果
SHOW INDEX FROM orders;
SHOW INDEX FROM factory_profiles; 