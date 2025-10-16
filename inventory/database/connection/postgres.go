package connection

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"time"

	"synapsis/inventory/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func connection(cfg config.AppConfig) (*gorm.DB, error) {
	var dbLogFile *os.File
	dbLogFile, err := os.OpenFile(fmt.Sprintf("%s/%s", cfg.Postgres.LogDirectory, cfg.Postgres.LogFile), os.O_CREATE|os.O_RDWR|os.O_APPEND, fs.ModePerm)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		_ = os.MkdirAll(cfg.Postgres.LogDirectory, os.ModePerm)
		dbLogFile, _ = os.Create(fmt.Sprintf("%s/%s", cfg.Postgres.LogDirectory, cfg.Postgres.LogFile))
		_ = dbLogFile.Chmod(fs.ModePerm)
	}

	dbLogger := logger.New(
		log.New(io.MultiWriter(os.Stdout, dbLogFile), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      false,
		},
	)

	loc, _ := time.LoadLocation("Asia/Jakarta")
	gormConfig := &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 dbLogger,
		NowFunc: func() time.Time {
			return time.Now().In(loc)
		},
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&timezone=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
		cfg.Postgres.Zone,
	)

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}

type Instance struct {
	GormDB *gorm.DB
}

type DBInstance interface {
	Database() *gorm.DB
	Close()
}

func NewDatabaseInstance(cfg config.AppConfig) DBInstance {
	ins := new(Instance)

	database, err := connection(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database. error: %s", err.Error()))
	}
	ins.GormDB = database

	return ins
}

func (i *Instance) Database() *gorm.DB {
	return i.GormDB
}

func (i *Instance) Close() {
	db, err := i.GormDB.DB()
	if err == nil {
		_ = db.Close()
	}
}
