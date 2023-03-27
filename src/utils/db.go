package utils

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbConnection *gorm.DB

const userNameEnvVar = "FAB_DB_USER_NAME"
const psswdEnvVar = "FAB_DB_PASSWORD"
const hostEnvVar = "FAB_DB_HOST_NAME"
const portEnvVar = "FAB_DB_PORT"

func InitDB() {
	envVarNames := [4]string{userNameEnvVar, psswdEnvVar, hostEnvVar, portEnvVar}
	envVarVals := ReadEnvVars(envVarNames[:])
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
		envVarVals[userNameEnvVar],
		envVarVals[psswdEnvVar],
		envVarVals[hostEnvVar],
		envVarVals[portEnvVar],
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	dbConnection = db
}

func GetDB() *gorm.DB {
	if dbConnection == nil {
		InitDB()
	}
	return dbConnection
}
