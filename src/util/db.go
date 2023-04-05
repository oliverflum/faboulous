package util

import (
	"fmt"

	"github.com/oliverflum/faboulous/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbConnection *gorm.DB

const userNameEnvVar = "FAB_DB_USER_NAME"
const psswdEnvVar = "FAB_DB_PASSWORD"
const hostEnvVar = "FAB_DB_HOST_NAME"
const portEnvVar = "FAB_DB_PORT"
const dbNameVar = "FAB_DB_DB_NAME"

func InitDB() *gorm.DB {
	envVarNames := [5]string{
		userNameEnvVar,
		psswdEnvVar,
		hostEnvVar,
		portEnvVar,
		dbNameVar,
	}
	envVarVals := ReadEnvVars(envVarNames[:])
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
	dbConnection = db
	dbConnection.AutoMigrate(&model.Feature{}, &model.Test{})
	return dbConnection
}

func GetDB() *gorm.DB {
	if dbConnection == nil {
		InitDB()
	}
	return dbConnection
}
