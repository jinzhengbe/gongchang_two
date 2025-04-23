-- 创建用户并授予权限
CREATE USER IF NOT EXISTS 'gongchang'@'%' IDENTIFIED BY 'gongchang';
GRANT ALL PRIVILEGES ON gongchang.* TO 'gongchang'@'%';
FLUSH PRIVILEGES; 