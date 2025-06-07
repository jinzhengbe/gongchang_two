-- 修改 users 表中 role 字段的默认值
ALTER TABLE users MODIFY COLUMN role varchar(191) NOT NULL; 