-- 为fabrics表添加designer_id字段
ALTER TABLE fabrics ADD COLUMN designer_id VARCHAR(191) NULL;
 
-- 添加索引以提高查询性能
CREATE INDEX idx_fabrics_designer_id ON fabrics(designer_id);
CREATE INDEX idx_fabrics_supplier_id ON fabrics(supplier_id); 