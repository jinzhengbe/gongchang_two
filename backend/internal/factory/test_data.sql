-- 创建订单表
CREATE TABLE IF NOT EXISTS orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(191) UNIQUE NOT NULL,
    designer_id BIGINT UNSIGNED NOT NULL,
    factory_id BIGINT UNSIGNED NOT NULL,
    status INT NOT NULL DEFAULT 1,
    total_amount DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_designer_id (designer_id),
    INDEX idx_factory_id (factory_id)
);

-- 创建订单项表
CREATE TABLE IF NOT EXISTS order_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL,
    product_name VARCHAR(191) NOT NULL,
    quantity INT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    INDEX idx_order_id (order_id)
);

-- 插入测试订单数据
INSERT INTO orders (
    title,
    description,
    fabric,
    quantity,
    factory_id,
    status,
    designer_id,
    unit_price,
    total_price,
    payment_status,
    order_type,
    order_date,
    created_at,
    updated_at
) VALUES
(
    '测试订单1',
    '这是一个测试订单',
    '棉布',
    2,
    '1',
    'processing',
    '1',
    500.00,
    1000.00,
    'unpaid',
    'standard',
    NOW(),
    NOW(),
    NOW()
),
(
    '测试订单2',
    '这是另一个测试订单',
    '丝绸',
    4,
    '1',
    'completed',
    '1',
    500.00,
    2000.00,
    'paid',
    'standard',
    NOW(),
    NOW(),
    NOW()
);

-- 插入测试订单项数据
INSERT INTO order_items (order_id, product_name, quantity, price) VALUES
(1, '测试产品1', 2, 500.00),
(2, '测试产品2', 4, 500.00); 