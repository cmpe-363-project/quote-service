package repository

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MysqlRepository struct {
	db *gorm.DB
}

// NewMysqlRepositoryConfig contains configuration for creating a new MySQL repository.
type NewMysqlRepositoryConfig struct {
	// Host is the MySQL server hostname or IP address.
	Host string

	// Port is the MySQL server port.
	Port int

	// Username is the MySQL user for authentication.
	Username string

	// Password is the MySQL user password.
	Password string

	// Database is the name of the database to connect to.
	Database string

	// MaxOpenConns sets the maximum number of open connections to the database.
	MaxOpenConns int

	// MaxIdleConns sets the maximum number of idle connections in the pool.
	MaxIdleConns int

	// ConnMaxLifetime sets the maximum amount of time a connection may be reused.
	ConnMaxLifetime time.Duration

	// LogLevel sets the GORM logger level. If nil, logger.Silent will be used.
	LogLevel *logger.LogLevel

	// AutoMigrate when true, automatically migrates the database schema.
	AutoMigrate bool
}

// NewMysqlRepository creates a new MySQL repository instance.
func NewMysqlRepository(config NewMysqlRepositoryConfig) (*MysqlRepository, error) {
	logLevel := logger.Silent
	if config.LogLevel != nil {
		logLevel = *config.LogLevel
	}

	// Build DSN with utf8mb4 charset and case-sensitive collation
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&collation=utf8mb4_bin&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{ //nolint:exhaustruct
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying database: %w", err)
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	repo := &MysqlRepository{db: db}

	if config.AutoMigrate {
		if err := AutoMigrate(db); err != nil {
			return nil, fmt.Errorf("failed to auto-migrate: %w", err)
		}
	}

	return repo, nil
}

// Close closes the database connection.
func (r *MysqlRepository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying database: %w", err)
	}

	return sqlDB.Close()
}

func (r *MysqlRepository) GetQuoteByID(id int) (*Quote, error) {
	var quote Quote
	if err := r.db.First(&quote, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &quote, nil
}

func (r *MysqlRepository) GetRandomQuote() (*Quote, error) {
	var quote Quote
	if err := r.db.Order("RAND()").First(&quote).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &quote, nil
}
