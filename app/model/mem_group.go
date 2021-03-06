package model

import (
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
)

const TableMemGroup = "mem_group"

type MemGroup struct {
	Id        int `gorm:"primary_key"`
	Code      string
	Name      string
	IsDefault int8
	IsEnable  int8
	CreatedAt carbon.DateTime
	UpdatedAt carbon.DateTime
	DeletedAt gorm.DeletedAt
}
