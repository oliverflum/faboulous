package util

import (
	"fmt"

	"github.com/oliverflum/faboulous/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbConnection *gorm.DB

const userNameEnvVar = "FAB_DB_USER_NAME"
const psswdEnvVar = "FAB_DB_PASSWORD"
const hostEnvVar = "FAB_DB_HOST_NAME"
const portEnvVar = "FAB_DB_PORT"
const dbNameVar = "FAB_DB_DB_NAME"
const dbTypeVar = "FAB_DB_TYPE"

func InitSqliteDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	return db
}

func InitMysqlDB() *gorm.DB {
	envVarNames := []string{
		userNameEnvVar,
		psswdEnvVar,
		hostEnvVar,
		portEnvVar,
		dbNameVar,
	}
	envVarVals := ReadEnvVars(envVarNames)
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		envVarVals[userNameEnvVar],
		envVarVals[psswdEnvVar],
		envVarVals[hostEnvVar],
		envVarVals[portEnvVar],
		envVarVals[dbNameVar],
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	return db
}

func InitDB() *gorm.DB {
	envTypeVarNameArr := []string{
		dbTypeVar,
	}
	dbType := ReadEnvVars(envTypeVarNameArr)[dbTypeVar]
	var db *gorm.DB
	if dbType == "sqlite" {
		db = InitSqliteDB()
	} else if dbType == "mysql" {
		db = InitMysqlDB()
	} else {
		panic("Invalid database type. Supported types are: sqlite, mysql")
	}
	db.AutoMigrate(&model.FeatureEntity{}, &model.TestEntity{}, &model.VariantEntity{}, &model.VariantFeatureEntity{})
	return db
}

func GetDB() *gorm.DB {
	if dbConnection == nil {
		dbConnection = InitDB()
	}
	return dbConnection
}
