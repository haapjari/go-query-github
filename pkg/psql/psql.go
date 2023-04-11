package psql

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
    "github.com/haapjari/go-query-github/pkg/models"
)

type PostgreSQL struct {
	Host           string
	Port           int
	User           string
	Password       string
	Database       string
	GormObject *gorm.DB
}

func NewPostgreSQL(host string, port int, user string, password string, database string) *PostgreSQL {
	p := &PostgreSQL{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
	}

	return p
}

func (p *PostgreSQL) Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Helsinki", p.Host, p.User, p.Password, p.Database, p.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}

func (p *PostgreSQL) Close() error {
    sqlDB, err := p.GormObject.DB()
    if err != nil {
        return err
    }

    return sqlDB.Close()
}

func (p *PostgreSQL) UpdateRows(db *gorm.DB, url string, column string, value int) error {
    result := db.Model(&models.Repo{}).Where("url = ?", url).Update(column, value)

    return result.Error
}
