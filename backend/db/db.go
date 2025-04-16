package db

import (
	"fmt"

	"github.com/oliverflum/faboulous/backend/model"
	"github.com/oliverflum/faboulous/backend/util"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbConnection *gorm.DB

const userNameEnvVar = "FAB_DB_USER_NAME"
const psswdEnvVar = "FAB_DB_PASSWORD"
const hostEnvVar = "FAB_DB_HOST_NAME"
const portEnvVar = "FAB_DB_PORT"
const dbNameVar = "FAB_DB_DB_NAME"
const dbTypeVar = "FAB_DB_TYPE"

func InitSqliteDB(inMemory bool, gormConfig *gorm.Config) *gorm.DB {
	var db *gorm.DB
	var err error

	if inMemory {
		db, err = gorm.Open(sqlite.Open(":memory:"), gormConfig)
	} else {
		db, err = gorm.Open(sqlite.Open("test.db"), gormConfig)
	}

	if err != nil {
		panic("Failed to connect database")
	}
	return db
}

func InitMysqlDB(gormConfig *gorm.Config) *gorm.DB {
	envVarNames := []string{
		userNameEnvVar,
		psswdEnvVar,
		hostEnvVar,
		portEnvVar,
		dbNameVar,
	}
	envVarVals := util.ReadEnvVars(envVarNames)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		envVarVals[userNameEnvVar],
		envVarVals[psswdEnvVar],
		envVarVals[hostEnvVar],
		envVarVals[portEnvVar],
		envVarVals[dbNameVar],
	)
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		panic("Failed to connect database")
	}
	return db
}

func InitDB() *gorm.DB {
	envTypeVarNameArr := []string{
		dbTypeVar,
	}
	dbType := util.ReadEnvVars(envTypeVarNameArr)[dbTypeVar]
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(util.GetLogLevelGorm()),
	}
	var db *gorm.DB
	if dbType == "sqlite" {
		db = InitSqliteDB(false, gormConfig)
	} else if dbType == "mysql" {
		db = InitMysqlDB(gormConfig)
	} else if dbType == "memory" {
		db = InitSqliteDB(true, gormConfig)
	} else {
		panic("Invalid database type. Supported types are: sqlite, mysql, memory")
	}
	db.AutoMigrate(&model.Feature{}, &model.Test{}, &model.Variant{}, &model.VariantFeature{})
	return db
}

func GetDB() *gorm.DB {
	if dbConnection == nil {
		dbConnection = InitDB()
	}
	return dbConnection
}
