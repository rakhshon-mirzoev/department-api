package models

import "time"

type Employee struct {
	ID           int64      `json:"id" gorm:"column:id;primaryKey"`
	DepartmentID int64      `json:"department_id" gorm:"column:department_id"`
	FullName     string     `json:"full_name" gorm:"column:full_name"`
	Position     string     `json:"position" gorm:"column:position"`
	HiredAt      *time.Time `json:"hired_at" gorm:"column:hired_at"`
	CreatedAt    time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}
