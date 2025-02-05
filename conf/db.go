package conf

import (
	"github.com/spf13/viper"
	mqysl "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func InitDb() (*gorm.DB, error) {
	LogMode := logger.Info
	if !viper.GetBool("mode.dev") {
		LogMode = logger.Error
	}
	mysqlInfo := mqysl.Open(viper.GetString("db.dsn"))
	db, err := gorm.Open(mysqlInfo, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "sys_",
			SingularTable: true,
		},
		Logger:                                   logger.Default.LogMode(LogMode),
		SkipDefaultTransaction:                   false,
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		return nil, err
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(viper.GetInt("mode.MaxIdleConn"))
	sqlDB.SetMaxOpenConns(viper.GetInt("mode.MaxOpenConn"))
	sqlDB.SetConnMaxLifetime(time.Hour)

	db.AutoMigrate(
	//&user_model.User{},
	//&mq_model.Mq{},
	//&item_model.Item{},
	//&item_model.Tag{},
	//&ssh_model.SSH{},
	//&project_model.Project{},
	//&module_model.Module{},
	//&case_model.Case{},
	//&case_step_model.CaseStep{},
	)

	return db, nil

}
