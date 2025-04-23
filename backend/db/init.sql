-- Create database if not exists
CREATE DATABASE IF NOT EXISTS sewingmast;
USE sewingmast;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE,
    role ENUM('admin', 'designer', 'factory', 'supplier') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create factories table
CREATE TABLE IF NOT EXISTS factories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    type ENUM('standard', 'premium') NOT NULL,
    location VARCHAR(255) NOT NULL,
    rating DECIMAL(2,1),
    certification TEXT,
    contact_info TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create orders table
CREATE TABLE IF NOT EXISTS orders (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    factory_id BIGINT NOT NULL,
    designer_id VARCHAR(64) NOT NULL,
    customer_id VARCHAR(64) NOT NULL,
    status ENUM('pending', 'accepted', 'in_progress', 'completed', 'cancelled') NOT NULL,
    details TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (factory_id) REFERENCES factories(id)
);

-- Create sample admin user
INSERT INTO users (username, password, email, role) 
VALUES ('admin', '$2a$10$xVB1H7qIJT.z8Xw9H3zOyeQZr0Kq8X9q9Y9X9q9Y9X9q9Y9X9q', 'admin@sewingmast.com', 'admin')
ON DUPLICATE KEY UPDATE username=username; 