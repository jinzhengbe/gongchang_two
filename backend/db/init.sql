-- Create database if not exists
CREATE DATABASE IF NOT EXISTS gongchang;
USE gongchang;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL,
    role ENUM('admin', 'designer', 'factory', 'supplier') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- Create factories table
CREATE TABLE IF NOT EXISTS factories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    type ENUM('standard', 'premium') NOT NULL,
    location VARCHAR(255) NOT NULL,
    rating DECIMAL(2,1),
    certification TEXT,
    contact_info TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

-- Create orders table
CREATE TABLE IF NOT EXISTS orders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    fabric VARCHAR(255),
    quantity INT NOT NULL,
    factory_id VARCHAR(36),
    designer_id VARCHAR(36) NOT NULL,
    customer_id VARCHAR(36) NOT NULL,
    unit_price DECIMAL(10,2),
    total_price DECIMAL(10,2),
    payment_status VARCHAR(50),
    shipping_address TEXT,
    order_type VARCHAR(50),
    fabrics TEXT,
    delivery_date TIMESTAMP NULL DEFAULT NULL,
    order_date TIMESTAMP NULL DEFAULT NULL,
    special_requirements TEXT,
    status ENUM('draft', 'published', 'completed', 'cancelled') NOT NULL DEFAULT 'draft',
    attachments JSON,
    models JSON,
    images JSON,
    videos JSON,
    created_at TIMESTAMP NULL DEFAULT NULL,
    updated_at TIMESTAMP NULL DEFAULT NULL,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (factory_id) REFERENCES users(id),
    FOREIGN KEY (designer_id) REFERENCES users(id),
    FOREIGN KEY (customer_id) REFERENCES users(id)
);

-- Create sample admin user
INSERT INTO users (id, username, password, email, role) 
VALUES (UUID(), 'admin', '$2a$10$xVB1H7qIJT.z8Xw9H3zOyeQZr0Kq8X9q9Y9X9q9Y9X9q9Y9X9q', 'admin@sewingmast.com', 'admin')
ON DUPLICATE KEY UPDATE username=username; 