package kitGorm

import (
	"gorm.io/gorm"
)

type ID struct {
	ID int64 `gorm:"column:id;primary_key" json:"id" form:"id"`
}

type BaseModel struct {
	ID        int64          `gorm:"column:id;primary_key" json:"id" form:"id"`
	CreatedAt Time           `gorm:"column:create_time" json:"create_time" form:"create_time"`
	UpdatedAt Time           `gorm:"column:update_time" json:"update_time" form:"update_time"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_time" sql:"index" json:"-"`
}

type DefaultModel struct {
	ID        int64          `gorm:"column:id;primary_key" json:"id" form:"id"`
	CreatedAt Time           `gorm:"column:created_at;default:null" json:"created_at" form:"created_at"`
	UpdatedAt Time           `gorm:"column:updated_at; default:null" json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `sql:"index"gorm:"column:deleted_at; default:null" json:"deleted_at" form:"deleted_at"`
}


type IdCompany struct {
	ID int64 `gorm:"column:id;primary_key" json:"id" form:"id"`
	CompanyId   int32 `gorm:"column:company_id" json:"company_id"`
}

