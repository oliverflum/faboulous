package model

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

type valueType string

const (
	BOOL   valueType = "BOOL"
	STRING valueType = "STRING"
	INT    valueType = "INT"
	FLOAT  valueType = "FLOAT"
)

func (self *valueType) Scan(value interface{}) error {
	*self = valueType(value.([]byte))
	return nil
}

func (self valueType) Value() (driver.Value, error) {
	return string(self), nil
}

type Feature struct {
	gorm.Model
	Name      string
	ValueType valueType `sql:"type:ENUM('BOOL', 'STRING', 'INT', 'FLOAT')" gorm:"column:car_type"`
}
