package database

import (
	"channel-service/internal/models"

	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	DB *gorm.DB
}

// NewPostgresDB accepts credentials as arguments
func NewPostgresDB(host, port, user, password, dbname string) *PostgresDB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	db.AutoMigrate(&models.Channel{}, &models.Subscription{})
	return &PostgresDB{DB: db}
}

func (p *PostgresDB) CreateChannel(ch *models.Channel) error {
	return p.DB.Create(ch).Error
}

func (p *PostgresDB) Close() {
	sqlDB, _ := p.DB.DB()
	sqlDB.Close()
}

func (p *PostgresDB) CreateSubscription(sub *models.Subscription) error {
	return p.DB.Create(sub).Error
}

func (p *PostgresDB) IncrementSubscribers(channelID string) error {
	return p.DB.Model(&models.Channel{}).
		Where("id = ?", channelID).
		Update("subscribers", gorm.Expr("subscribers + 1")).Error
}
