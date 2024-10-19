package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbService struct {
	db *gorm.DB
}

type PostresConf struct {
	User string
	Pass string
	Host string
	Port int
	Name string
}

func New(conf PostresConf) (*DbService, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", conf.Host, conf.User, conf.Pass, conf.Name, conf.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	dbService := DbService{
		db: db,
	}

	db.AutoMigrate(&InverterReading{})
	db.AutoMigrate(&HeaterLogs{})

	return &dbService, nil
}

func (s DbService) InsertReading(reading int) error {
	record := InverterReading{
		Reading:   reading,
		CreatedAt: time.Now(),
	}

	s.db.Create(&record)

	return nil
}

func (s DbService) InsertHeaterlog(reading int, on bool) error {
	status := POWER_OFF
	if on {
		status = POWER_ON
	}
	record := HeaterLogs{
		Reading:   reading,
		CreatedAt: time.Now(),
		Status:    status,
	}

	s.db.Create(&record)

	return nil
}

func (s DbService) GetLastHeaterStatus() HeaterStatus {
	var log HeaterLogs
	s.db.Order("CreatedAt desc").First(&log)
	return log.Status
}
