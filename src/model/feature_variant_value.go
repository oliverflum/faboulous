package model

import "gorm.io/gorm"

type VariantFeatureEntity struct {
	gorm.Model
	FeatureID uint
	Type      string
	Value     any
}

func (variantFeature *VariantFeatureEntity) TableName() string {
	return "variant_features"
}
