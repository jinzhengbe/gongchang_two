package database

import (
	"context"
	"fmt"
	"log"
	"time"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseStats 数据库统计信息
type DatabaseStats struct {
	MaxOpenConnections int
	OpenConnections    int
	InUseConnections   int
	IdleConnections    int
	WaitCount          int64
	WaitDuration       time.Duration
	MaxIdleClosed      int64
	MaxLifetimeClosed  int64
}

// MonitorDatabase 监控数据库性能
func MonitorDatabase(db *gorm.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		stats, err := GetDatabaseStats(db)
		if err != nil {
			log.Printf("Failed to get database stats: %v", err)
			continue
		}

		log.Printf("Database Stats:")
		log.Printf("  Max Open Connections: %d", stats.MaxOpenConnections)
		log.Printf("  Open Connections: %d", stats.OpenConnections)
		log.Printf("  In Use Connections: %d", stats.InUseConnections)
		log.Printf("  Idle Connections: %d", stats.IdleConnections)
		log.Printf("  Wait Count: %d", stats.WaitCount)
		log.Printf("  Wait Duration: %v", stats.WaitDuration)
		log.Printf("  Max Idle Closed: %d", stats.MaxIdleClosed)
		log.Printf("  Max Lifetime Closed: %d", stats.MaxLifetimeClosed)
	}
}

// GetDatabaseStats 获取数据库统计信息
func GetDatabaseStats(db *gorm.DB) (*DatabaseStats, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}

	stats := &DatabaseStats{}
	stats.MaxOpenConnections = sqlDB.Stats().MaxOpenConnections
	stats.OpenConnections = sqlDB.Stats().OpenConnections
	stats.InUseConnections = sqlDB.Stats().InUse
	stats.IdleConnections = sqlDB.Stats().Idle
	stats.WaitCount = sqlDB.Stats().WaitCount
	stats.WaitDuration = sqlDB.Stats().WaitDuration
	stats.MaxIdleClosed = sqlDB.Stats().MaxIdleClosed
	stats.MaxLifetimeClosed = sqlDB.Stats().MaxLifetimeClosed

	return stats, nil
}

// MonitorSlowQueries 监控慢查询
func MonitorSlowQueries(db *gorm.DB, threshold time.Duration) {
	// 设置慢查询阈值
	db = db.Session(&gorm.Session{
		Logger: logger.Default.LogMode(logger.Info),
	})

	// 监控查询执行时间
	db.Callback().Query().Before("gorm:query").Register("monitor:before_query", func(db *gorm.DB) {
		db.Statement.Context = context.WithValue(db.Statement.Context, "start_time", time.Now())
	})

	db.Callback().Query().After("gorm:query").Register("monitor:after_query", func(db *gorm.DB) {
		startTime := db.Statement.Context.Value("start_time").(time.Time)
		duration := time.Since(startTime)

		if duration > threshold {
			log.Printf("Slow Query detected (%.2fs): %s", duration.Seconds(), db.Statement.SQL.String())
		}
	})
} 