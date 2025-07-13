-- 修复工厂权限问题的SQL脚本
-- 问题：工厂记录的id和user_id字段不一致，导致权限校验失败

-- 1. 查看当前工厂记录
SELECT id, user_id, company_name FROM users WHERE role = 'factory';

-- 2. 修复工厂记录，确保id和user_id一致
-- 对于gongchang用户，将工厂记录的id设置为user_id
UPDATE users 
SET id = user_id 
WHERE username = 'gongchang' AND role = 'factory';

-- 3. 验证修复结果
SELECT id, user_id, company_name, username FROM users WHERE role = 'factory';

-- 4. 如果上述方法不行，可以尝试另一种方法：修改user_id为id
-- UPDATE users 
-- SET user_id = id 
-- WHERE username = 'gongchang' AND role = 'factory';

-- 5. 检查是否有重复的user_id
SELECT user_id, COUNT(*) as count 
FROM users 
WHERE role = 'factory' 
GROUP BY user_id 
HAVING count > 1;

-- 6. 如果需要，可以重新创建工厂记录
-- DELETE FROM users WHERE username = 'gongchang' AND role = 'factory';
-- INSERT INTO users (id, user_id, username, email, role, company_name, address, capacity, equipment, certificates, photos, videos, employee_count, rating, status, created_at, updated_at) 
-- VALUES ('3af8e32a-e267-45f1-8959-faf3f0787bfa', '3af8e32a-e267-45f1-8959-faf3f0787bfa', 'gongchang', '010759839@factory.com', 'factory', '工场部分更新测试', '广东省深圳市南山区科技园', 3000, '全自动裁剪机,工业缝纫机,质检设备,激光切割机,智能熨烫设备', 'ISO9001,质量管理体系认证,环保认证,OHSAS18001,CE认证', '[]', '[]', 300, 0, 1, NOW(), NOW()); 