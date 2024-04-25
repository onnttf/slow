package model

import (
	"time"
)

const TableNamePetrol = "petrol"

// Petrol mapped from table <petrol>
type Petrol struct {
	ID          int32     `gorm:"column:id;primaryKey;autoIncrement:true;comment:unique identifier" json:"id"`                           // unique identifier
	Type        string    `gorm:"column:type;not null;index:idx_type_area,priority:1;comment:petrol type" json:"type"`                   // petrol type
	Area        string    `gorm:"column:area;not null;index:idx_type_area,priority:2;comment:area" json:"area"`                          // area
	ReleaseDate string    `gorm:"column:release_date;not null;comment:release_date" json:"release_date"`                                 // release_date
	Price       string    `gorm:"column:price;not null;default:0.00;comment:price" json:"price"`                                         // price
	CreatedAt   time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;comment:created timestamp" json:"created_at"`      // created timestamp
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP;comment:last updated timestamp" json:"updated_at"` // last updated timestamp
}

// TableName Petrol's table name
func (*Petrol) TableName() string {
	return TableNamePetrol
}
