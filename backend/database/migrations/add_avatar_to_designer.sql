-- 为设计师表添加头像字段
-- 迁移脚本：添加头像字段到 designer_profiles 表

-- 检查字段是否已存在，如果不存在则添加
SET @sql = (SELECT IF(
    (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS 
     WHERE TABLE_SCHEMA = 'gongchang' 
     AND TABLE_NAME = 'designer_profiles' 
     AND COLUMN_NAME = 'avatar') = 0,
    'ALTER TABLE designer_profiles ADD COLUMN avatar VARCHAR(500) DEFAULT "" COMMENT "头像URL"',
    'SELECT "avatar字段已存在"'
));

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 验证字段是否添加成功
SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE, COLUMN_DEFAULT, COLUMN_COMMENT
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_SCHEMA = 'gongchang' 
AND TABLE_NAME = 'designer_profiles' 
AND COLUMN_NAME = 'avatar'; 