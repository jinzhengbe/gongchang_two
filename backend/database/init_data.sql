-- 布料分类数据
INSERT INTO fabric_categories (name, description, icon, sort, status, created_at, updated_at) VALUES
('棉布', '天然棉纤维制成的布料，透气性好，适合制作夏季服装', 'cotton', 1, 1, NOW(), NOW()),
('丝绸', '天然蚕丝制成的布料，质地柔软，光泽度高', 'silk', 2, 1, NOW(), NOW()),
('羊毛', '天然羊毛纤维制成的布料，保暖性好', 'wool', 3, 1, NOW(), NOW()),
('麻布', '天然麻纤维制成的布料，透气性好，适合夏季', 'linen', 4, 1, NOW(), NOW()),
('化纤', '化学纤维制成的布料，价格便宜，易打理', 'synthetic', 5, 1, NOW(), NOW()),
('混纺', '天然纤维与化学纤维混合制成的布料', 'blend', 6, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();

-- 布料数据
INSERT INTO fabrics (name, category, material, color, pattern, weight, width, price, unit, stock, min_order, description, image_url, thumbnail_url, tags, status, created_at, updated_at) VALUES
('纯棉平纹布', '棉布', '100%棉', '白色', '平纹', 120.00, 150.00, 15.50, '米', 100, 1, '优质纯棉平纹布，透气性好，适合制作衬衫、T恤等', '/uploads/fabrics/cotton_plain_white.jpg', '/uploads/fabrics/thumbnails/cotton_plain_white.jpg', '棉布,白色,平纹,透气', 1, NOW(), NOW()),
('纯棉斜纹布', '棉布', '100%棉', '蓝色', '斜纹', 180.00, 150.00, 22.00, '米', 80, 1, '纯棉斜纹布，质地厚实，适合制作牛仔裤、工装裤等', '/uploads/fabrics/cotton_twill_blue.jpg', '/uploads/fabrics/thumbnails/cotton_twill_blue.jpg', '棉布,蓝色,斜纹,厚实', 1, NOW(), NOW()),
('真丝缎面', '丝绸', '100%桑蚕丝', '红色', '缎面', 16.00, 140.00, 85.00, '米', 50, 1, '高档真丝缎面，光泽度好，适合制作礼服、旗袍等', '/uploads/fabrics/silk_satin_red.jpg', '/uploads/fabrics/thumbnails/silk_satin_red.jpg', '丝绸,红色,缎面,高档', 1, NOW(), NOW()),
('真丝雪纺', '丝绸', '100%桑蚕丝', '粉色', '雪纺', 12.00, 140.00, 65.00, '米', 60, 1, '轻薄真丝雪纺，飘逸感强，适合制作连衣裙、衬衫等', '/uploads/fabrics/silk_chiffon_pink.jpg', '/uploads/fabrics/thumbnails/silk_chiffon_pink.jpg', '丝绸,粉色,雪纺,轻薄', 1, NOW(), NOW()),
('羊毛呢', '羊毛', '100%羊毛', '灰色', '呢面', 280.00, 150.00, 120.00, '米', 40, 1, '优质羊毛呢，保暖性好，适合制作大衣、西装等', '/uploads/fabrics/wool_grey.jpg', '/uploads/fabrics/thumbnails/wool_grey.jpg', '羊毛,灰色,呢面,保暖', 1, NOW(), NOW()),
('亚麻布', '麻布', '100%亚麻', '米色', '平纹', 200.00, 150.00, 35.00, '米', 70, 1, '天然亚麻布，透气性好，适合制作夏季服装', '/uploads/fabrics/linen_beige.jpg', '/uploads/fabrics/thumbnails/linen_beige.jpg', '麻布,米色,平纹,透气', 1, NOW(), NOW()),
('涤纶面料', '化纤', '100%涤纶', '黑色', '平纹', 150.00, 150.00, 12.00, '米', 120, 1, '涤纶面料，价格便宜，易打理，适合制作工作服等', '/uploads/fabrics/polyester_black.jpg', '/uploads/fabrics/thumbnails/polyester_black.jpg', '化纤,黑色,平纹,便宜', 1, NOW(), NOW()),
('棉麻混纺', '混纺', '55%棉+45%麻', '绿色', '平纹', 160.00, 150.00, 28.00, '米', 90, 1, '棉麻混纺布，结合了棉的柔软和麻的透气性', '/uploads/fabrics/cotton_linen_green.jpg', '/uploads/fabrics/thumbnails/cotton_linen_green.jpg', '混纺,绿色,平纹,透气', 1, NOW(), NOW()),
('丝绸印花', '丝绸', '100%桑蚕丝', '花色', '印花', 14.00, 140.00, 75.00, '米', 45, 1, '真丝印花面料，图案精美，适合制作连衣裙、衬衫等', '/uploads/fabrics/silk_print_colorful.jpg', '/uploads/fabrics/thumbnails/silk_print_colorful.jpg', '丝绸,花色,印花,精美', 1, NOW(), NOW()),
('牛仔布', '棉布', '98%棉+2%氨纶', '蓝色', '斜纹', 250.00, 150.00, 25.00, '米', 85, 1, '弹力牛仔布，适合制作牛仔裤、夹克等', '/uploads/fabrics/denim_blue.jpg', '/uploads/fabrics/thumbnails/denim_blue.jpg', '棉布,蓝色,斜纹,弹力', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW(); 