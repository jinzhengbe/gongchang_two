-- 工厂搜索相关数据库索引
-- 用于优化工厂搜索API的性能

-- 1. 工厂搜索索引
-- 为工厂名称和地址创建复合索引，支持关键词搜索
CREATE INDEX IF NOT EXISTS idx_factory_profiles_search 
ON factory_profiles (company_name, address);

-- 2. 工厂专业领域索引
-- 为工厂专业领域表创建索引，支持专业领域筛选
CREATE INDEX IF NOT EXISTS idx_factory_specialties_factory_id 
ON factory_specialties (factory_id);

CREATE INDEX IF NOT EXISTS idx_factory_specialties_specialty 
ON factory_specialties (specialty);

-- 3. 工厂地区索引
-- 为工厂地址创建索引，支持地区筛选
CREATE INDEX IF NOT EXISTS idx_factory_profiles_address 
ON factory_profiles (address);

-- 4. 工厂评分索引
-- 为工厂评分表创建索引，支持评分筛选和排序
CREATE INDEX IF NOT EXISTS idx_factory_ratings_factory_id 
ON factory_ratings (factory_id);

CREATE INDEX IF NOT EXISTS idx_factory_ratings_rating 
ON factory_ratings (rating);

-- 5. 工厂状态索引
-- 为工厂状态创建索引，支持合作状态筛选
CREATE INDEX IF NOT EXISTS idx_factory_profiles_status 
ON factory_profiles (status);

-- 6. 工厂创建时间索引
-- 为工厂创建时间创建索引，支持时间排序
CREATE INDEX IF NOT EXISTS idx_factory_profiles_created_at 
ON factory_profiles (created_at);

-- 7. 用户角色索引
-- 为用户角色创建索引，确保只搜索工厂用户
CREATE INDEX IF NOT EXISTS idx_users_role 
ON users (role);

-- 8. 复合索引优化
-- 为常用的组合查询创建复合索引
CREATE INDEX IF NOT EXISTS idx_factory_profiles_status_rating 
ON factory_profiles (status, rating);

CREATE INDEX IF NOT EXISTS idx_factory_profiles_address_status 
ON factory_profiles (address, status);

-- 9. 全文搜索索引（如果数据库支持）
-- 为工厂名称和地址创建全文搜索索引
-- 注意：这需要数据库支持全文搜索功能
-- CREATE FULLTEXT INDEX IF NOT EXISTS idx_factory_profiles_fulltext 
-- ON factory_profiles (company_name, address);

-- 10. 评分统计索引
-- 为评分统计查询创建索引
CREATE INDEX IF NOT EXISTS idx_factory_ratings_factory_rating 
ON factory_ratings (factory_id, rating);

-- 11. 专业领域复合索引
-- 为专业领域筛选创建复合索引
CREATE INDEX IF NOT EXISTS idx_factory_specialties_factory_specialty 
ON factory_specialties (factory_id, specialty);

-- 12. 用户关联索引
-- 为工厂与用户的关联查询创建索引
CREATE INDEX IF NOT EXISTS idx_factory_profiles_user_id 
ON factory_profiles (user_id);

-- 13. 删除标记索引
-- 为软删除查询创建索引
CREATE INDEX IF NOT EXISTS idx_users_deleted_at 
ON users (deleted_at);

-- 14. 更新时间索引
-- 为按更新时间排序创建索引
CREATE INDEX IF NOT EXISTS idx_factory_profiles_updated_at 
ON factory_profiles (updated_at);

-- 15. 容量索引
-- 为产能筛选创建索引
CREATE INDEX IF NOT EXISTS idx_factory_profiles_capacity 
ON factory_profiles (capacity);

-- 注释说明：
-- 这些索引将显著提升工厂搜索API的性能
-- 特别是在以下场景：
-- 1. 关键词搜索（工厂名称、地址）
-- 2. 地区筛选
-- 3. 专业领域筛选
-- 4. 评分筛选和排序
-- 5. 合作状态筛选
-- 6. 分页查询
-- 7. 排序操作

-- 使用建议：
-- 1. 在生产环境部署前，请测试这些索引对现有查询的影响
-- 2. 根据实际查询模式，可能需要调整或删除某些索引
-- 3. 定期监控索引使用情况，删除未使用的索引
-- 4. 考虑数据库的存储空间和写入性能影响 