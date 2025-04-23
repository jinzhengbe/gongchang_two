#!/bin/bash

# 设置备份目录
BACKUP_DIR="/backup/$(date +%Y-%m-%d)"
mkdir -p $BACKUP_DIR

# 备份数据库
echo "Backing up database..."
docker exec gongchang_mysql_1 mysqldump -u root -p123456 gongchang > $BACKUP_DIR/database.sql

# 备份日志文件
echo "Backing up logs..."
docker run --rm -v gongchang_backend_logs:/source -v $BACKUP_DIR:/backup alpine tar -czf /backup/backend_logs.tar.gz -C /source .
docker run --rm -v gongchang_mysql_logs:/source -v $BACKUP_DIR:/backup alpine tar -czf /backup/mysql_logs.tar.gz -C /source .

# 备份配置文件
echo "Backing up config files..."
tar -czf $BACKUP_DIR/config.tar.gz \
    backend/.env \
    docker-compose.yml \
    backend/config/config.yaml

# 设置权限
chmod -R 755 $BACKUP_DIR

# 清理旧备份（保留最近30天）
find /backup -type d -mtime +30 -exec rm -rf {} \;

echo "Backup completed: $BACKUP_DIR" 