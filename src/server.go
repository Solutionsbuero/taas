package ttrn

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Run runs the server application.
func Run(cfg Config, doDebug bool) {
	if doDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	turnoutPositionEvents := make(chan TurnoutPositionEvent)
	trainSpeedEvents := make(chan TrainSpeedEvent)
	trainPositionEvents := make(chan TrainPositionEvent)

	web := NewWeb(cfg, doDebug, turnoutPositionEvents, trainSpeedEvents, trainPositionEvents)
	mqtt := NewMqtt(cfg, turnoutPositionEvents, trainSpeedEvents, trainPositionEvents)

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
