package model

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

type Method string

const (
	HASH   Method = "HASH"
	POOL   Method = "POOL"
	RANDOM Method = "RANDOM"
)

func (self *Method) Scan(method interface{}) error {
	*self = Method(method.([]byte))
	return nil
}

func (self Method) Value() (driver.Value, error) {
	return string(self), nil
}

type Test struct {
	gorm.Model
	Active   bool
	Method   Method `gorm:"varchar(10)" validate:"required, oneof HASH POOL RANDOM"`
	Variants []VariantEntity
}
