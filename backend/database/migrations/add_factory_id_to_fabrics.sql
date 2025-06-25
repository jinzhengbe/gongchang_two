-- 为fabrics表添加factory_id字段
ALTER TABLE fabrics ADD COLUMN factory_id VARCHAR(191) NULL;
 
-- 添加索引以提高查询性能
CREATE INDEX idx_fabrics_factory_id ON fabrics(factory_id); 