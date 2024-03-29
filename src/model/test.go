package model

import (
	"gorm.io/gorm"
)

type Test struct {
	gorm.Model
	Size   int  `validate:"required,min=1,max=50" json:"name"`
	Active bool `json:"active" validate:"required"`
}
