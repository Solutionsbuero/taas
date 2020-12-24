package ttrn

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Run runs the server application.
func Run(cfg Config, doDebug bool) {
	turnoutEvents := make(chan TurnoutEvent)
	trainEvents := make(chan TrainEvent)

	web := NewWeb(cfg, doDebug, openDb(cfg), turnoutEvents, trainEvents)
	mqtt := NewMqtt(cfg, turnoutEvents, trainEvents)

	go mqtt.Run()
	web.Run()
}

func openDb(cfg Config) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(cfg.Db), &gorm.Config{})
	if err != nil {
		logrus.Panicf("failed to connect to db, %s", err)
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&State{})
	return db
}
