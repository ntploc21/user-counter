package app

import (
	"context"
	"fmt"
	"locntp-user-counter/config"
	"locntp-user-counter/internal/models"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMySqlDatabase() (*gorm.DB, error) {
	cfg := config.GetAppConfig().Database

	// Build DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DatabaseName,
	)

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		logrus.WithError(err).Error("Failed to connect to MySQL database")
		return nil, err
	}

	// Get generic database object to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		logrus.WithError(err).Error("Failed to get database instance")
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto migrate the schema
	if err := db.AutoMigrate(&models.User{}); err != nil {
		logrus.WithError(err).Error("Failed to auto migrate database schema")
		return nil, err
	}

	logrus.Info("Successfully connected to MySQL database")
	return db, nil
}

func CloseMySqlDatabase(database *gorm.DB) error {
	if database == nil {
		return nil
	}

	sqlDB, err := database.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Close(); err != nil {
		logrus.WithError(err).Error("Failed to close MySQL database")
		return err
	}

	logrus.Info("MySQL database connection closed")
	return nil
}

func NewRedisInMemoryDatabase() (*redis.Client, *redis_rate.Limiter, error) {
	cfg := config.GetAppConfig().Redis

	// Support both single and cluster mode
	var client *redis.Client
	var limiter *redis_rate.Limiter
	if len(cfg.Addrs) == 1 {
		// Single node mode
		client = redis.NewClient(&redis.Options{
			Addr:         cfg.Addrs[0],
			Password:     cfg.Password,
			DB:           0,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolSize:     10,
			MinIdleConns: 5,
		})
	} else {
		// For simplicity, use the first address
		client = redis.NewClient(&redis.Options{
			Addr:         cfg.Addrs[0],
			Password:     cfg.Password,
			DB:           0,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolSize:     10,
			MinIdleConns: 5,
		})
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logrus.WithError(err).Error("Failed to connect to Redis")
		return nil, nil, err
	}

	limiter = redis_rate.NewLimiter(client)

	logrus.Info("Successfully connected to Redis")
	return client, limiter, nil
}

func CloseRedisInMemoryDatabase(redis *redis.Client) error {
	if redis == nil {
		return nil
	}

	if err := redis.Close(); err != nil {
		logrus.WithError(err).Error("Failed to close Redis connection")
		return err
	}

	logrus.Info("Redis connection closed")
	return nil
}

func CreateIndexes(db *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.WithContext(ctx).AutoMigrate(&models.User{}); err != nil {
		return err
	}

	return nil
}
