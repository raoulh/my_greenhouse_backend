package models

import (
	"git.raoulh.pw/raoulh/my_greenhouse_backend/config"
	logger "git.raoulh.pw/raoulh/my_greenhouse_backend/log"

	gormlogger "gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

var (
	db      *gorm.DB
	logging *logrus.Entry
)

func init() {
	logging = logger.NewLogger("database")
}

//Init models
func Init() (err error) {
	logging.Infof("Opening database %s", config.Config.String("database.dsn"))
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN: config.Config.String("database.dsn"),
	}), &gorm.Config{
		//Logger: logger.NewGorm(), //TODO: our logrus impl does not work
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		return
	}

	migrateDb()

	return
}

func migrateDb() {
	//Migrate all tables
	db.AutoMigrate(
		&AuthToken{},
	)

	logging.Infof("Migration did run successfully")
}
