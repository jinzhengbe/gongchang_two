-- 设计师搜索相关数据库索引
-- 用于优化设计师搜索API的性能

-- 1. 设计师搜索索引
-- 为设计师名称和地址创建复合索引，支持关键词搜索
CREATE INDEX IF NOT EXISTS idx_designer_profiles_search 
ON designer_profiles (company_name, address);

-- 2. 设计师专业领域索引
-- 为设计师专业领域表创建索引，支持专业领域筛选
CREATE INDEX IF NOT EXISTS idx_designer_specialties_designer_id 
ON designer_specialties (designer_id);

CREATE INDEX IF NOT EXISTS idx_designer_specialties_specialty 
ON designer_specialties (specialty);

-- 3. 设计师地区索引
-- 为设计师地址创建索引，支持地区筛选
CREATE INDEX IF NOT EXISTS idx_designer_profiles_address 
ON designer_profiles (address);

-- 4. 设计师评分索引
-- 为设计师评分表创建索引，支持评分筛选和排序
CREATE INDEX IF NOT EXISTS idx_designer_ratings_designer_id 
ON designer_ratings (designer_id);

CREATE INDEX IF NOT EXISTS idx_designer_ratings_rating 
ON designer_ratings (rating);

-- 5. 设计师状态索引
-- 为设计师状态创建索引，支持状态筛选
CREATE INDEX IF NOT EXISTS idx_designer_profiles_status 
ON designer_profiles (status);

-- 6. 设计师创建时间索引
-- 为设计师创建时间创建索引，支持时间排序
CREATE INDEX IF NOT EXISTS idx_designer_profiles_created_at 
ON designer_profiles (created_at);

-- 7. 用户角色索引
-- 为用户角色创建索引，确保只搜索设计师用户
CREATE INDEX IF NOT EXISTS idx_users_role_designer 
ON users (role);

-- 8. 复合索引优化
-- 为常用的组合查询创建复合索引
CREATE INDEX IF NOT EXISTS idx_designer_profiles_status_rating 
ON designer_profiles (status, rating);

CREATE INDEX IF NOT EXISTS idx_designer_profiles_address_status 
ON designer_profiles (address, status);

-- 9. 评分统计索引
-- 为评分统计查询创建索引
CREATE INDEX IF NOT EXISTS idx_designer_ratings_designer_rating 
ON designer_ratings (designer_id, rating);

-- 10. 专业领域复合索引
-- 为专业领域筛选创建复合索引
CREATE INDEX IF NOT EXISTS idx_designer_specialties_designer_specialty 
ON designer_specialties (designer_id, specialty);

-- 11. 用户关联索引
-- 为设计师与用户的关联查询创建索引
CREATE INDEX IF NOT EXISTS idx_designer_profiles_user_id 
ON designer_profiles (user_id);

-- 12. 删除标记索引
-- 为软删除查询创建索引
CREATE INDEX IF NOT EXISTS idx_users_deleted_at_designer 
ON users (deleted_at);

-- 13. 更新时间索引
-- 为按更新时间排序创建索引
CREATE INDEX IF NOT EXISTS idx_designer_profiles_updated_at 
ON designer_profiles (updated_at);

-- 14. 评分者索引
-- 为评分者查询创建索引
CREATE INDEX IF NOT EXISTS idx_designer_ratings_rater_id 
ON designer_ratings (rater_id);

-- 15. 评分时间索引
-- 为按评分时间排序创建索引
CREATE INDEX IF NOT EXISTS idx_designer_ratings_created_at 
ON designer_ratings (created_at);

-- 注释说明：
-- 这些索引将显著提升设计师搜索API的性能
-- 特别是在以下场景：
-- 1. 关键词搜索（设计师名称、地址）
-- 2. 地区筛选
-- 3. 专业领域筛选
-- 4. 评分筛选和排序
-- 5. 状态筛选
-- 6. 分页查询
-- 7. 排序操作

-- 使用建议：
-- 1. 在生产环境部署前，请测试这些索引对现有查询的影响
-- 2. 根据实际查询模式，可能需要调整或删除某些索引
-- 3. 定期监控索引使用情况，删除未使用的索引
-- 4. 考虑数据库的存储空间和写入性能影响 