package database

import (
	"backend/config"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// BackupDatabase 执行数据库备份
func BackupDatabase(cfg *config.Config) error {
	// 创建备份目录
	backupDir := filepath.Join("/runData/gongChang/backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("backup_%s.sql", timestamp))

	// 构建 mysqldump 命令
	cmd := exec.Command("mysqldump",
		"-h", cfg.Database.Host,
		"-P", cfg.Database.Port,
		"-u", cfg.Database.User,
		"-p"+cfg.Database.Password,
		cfg.Database.DBName,
		"--single-transaction",
		"--quick",
		"--lock-tables=false",
		"--routines",
		"--triggers",
		"--events",
		"--add-drop-database",
		"--databases",
		"--result-file="+backupFile,
	)

	// 执行备份
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to backup database: %v", err)
	}

	return nil
}

// RestoreDatabase 从备份文件恢复数据库
func RestoreDatabase(cfg *config.Config, backupFile string) error {
	// 检查备份文件是否存在
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupFile)
	}

	// 构建 mysql 恢复命令
	cmd := exec.Command("mysql",
		"-h", cfg.Database.Host,
		"-P", cfg.Database.Port,
		"-u", cfg.Database.User,
		"-p"+cfg.Database.Password,
		cfg.Database.DBName,
		"-e", fmt.Sprintf("source %s", backupFile),
	)

	// 执行恢复
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restore database: %v", err)
	}

	return nil
}

// ListBackups 列出所有可用的备份
func ListBackups() ([]string, error) {
	backupDir := filepath.Join("/runData/gongChang/backups")
	files, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %v", err)
	}

	var backups []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			backups = append(backups, filepath.Join(backupDir, file.Name()))
		}
	}

	return backups, nil
}

// CleanOldBackups 清理旧的备份文件
func CleanOldBackups(maxAge time.Duration) error {
	backupDir := filepath.Join("/runData/gongChang/backups")
	files, err := os.ReadDir(backupDir)
	if err != nil {
		return fmt.Errorf("failed to read backup directory: %v", err)
	}

	now := time.Now()
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			info, err := file.Info()
			if err != nil {
				continue
			}

			if now.Sub(info.ModTime()) > maxAge {
				backupFile := filepath.Join(backupDir, file.Name())
				if err := os.Remove(backupFile); err != nil {
					return fmt.Errorf("failed to remove old backup %s: %v", backupFile, err)
				}
			}
		}
	}

	return nil
} 