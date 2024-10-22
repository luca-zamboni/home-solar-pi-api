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

	var err error
	var db *gorm.DB
	retry := 0
	for err != nil || retry == 0 {

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

		retry++

		if retry > 5 {
			return nil, err
		}

		time.Sleep(time.Second * 5)

	}

	if err != nil {
		return nil, err
	}

	dbService := DbService{
		db: db,
	}

	db.AutoMigrate(&InverterReading{})
	db.AutoMigrate(&HeaterAction{})
	db.AutoMigrate(&Action{})

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

func (s DbService) InsertHeaterAction(reading int, on bool) error {
	status := POWER_OFF
	if on {
		status = POWER_ON
	}

	now := time.Now()
	action := Action{
		CreatedAt: now,
		Name:      string(status),
	}

	heaterAction := HeaterAction{
		Reading:   reading,
		CreatedAt: now,
		Status:    status,
		Action:    action,
	}

	s.db.Create(&action)
	s.db.Create(&heaterAction)

	return nil
}
