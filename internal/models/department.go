package models

import "time"

type Department struct {
	ID        int64     `json:"id" gorm:"column:id;primaryKey"`
	Name      string    `json:"name" gorm:"column:name"`
	ParentID  *int64    `json:"parent_id" gorm:"column:parent_id"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

type DepartmentResponse struct {
	Department
	Employees []Employee           `json:"employees,omitempty"`
	Children  []DepartmentResponse `json:"children"`
}

type UpdateDepartment struct {
	Name     *string `json:"name"`
	ParentID *int64  `json:"parent_id"`
}
