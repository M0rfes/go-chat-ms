package db

import (
	"github.com/M0rfes/go-chat-ms/pkg/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB interface {
	Connect() error
	Close() error
	CreteMessage(*models.Message) error
}

type db struct {
	dns string
	db  *gorm.DB
}

func NewDB(dns string) DB {
	return &db{
		dns: dns,
	}
}

func (d *db) Connect() error {
	db, err := gorm.Open(postgres.Open(d.dns), &gorm.Config{})
	if err != nil {
		return err
	}
	d.db = db
	return nil
}

func (d *db) Close() error {
	db, err := d.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (d *db) CreteMessage(msg *models.Message) error {
	return d.db.Create(msg).Error
}
