-- 更新工厂Photos字段，添加测试图片
UPDATE factory_profiles 
SET photos = '["/uploads/00b6eebd-d058-4086-a715-13491e66d9ee.png", "/uploads/016efa37-8d49-4e31-9e78-1c8ddfae0b15.png", "/uploads/01995f2f-71f6-442f-85d8-ad96f32a632d.jpg", "/uploads/04f9ddf4-ef5d-405e-bd80-d80d2b23a6f1.jpg", "/uploads/063f5d78-c1b3-4cd9-82a5-d23458d36da5.jpg"]'
WHERE id = 1;

-- 验证更新结果
SELECT id, company_name, photos FROM factory_profiles WHERE id = 1; 