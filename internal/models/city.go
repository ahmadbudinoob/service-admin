package models

import (
	"time"
)

// City represents the TLD_MST_CITY table
type City struct {
	CityCD         int       `json:"city_code" db:"CITY_CD"`
	CityName       string    `json:"city_name" db:"CITY_NAME"`
	ProvinceCode   string    `json:"province_code" db:"PROVINCE_CD"`
	ProvinceName   string    `json:"province_name" db:"PROVINCE_NAME"`
	CreatedDate    time.Time `json:"created_date" db:"CRE_DT"`
	UserID         string    `json:"user_id" db:"USER_ID"`
	UpdatedDate    time.Time `json:"updated_date" db:"UPD_DT"`
	UpdatedBy      string    `json:"updated_by" db:"UPD_BY"`
	ApprovedDate   time.Time `json:"approved_date" db:"APPROVED_DT"`
	ApprovedBy     string    `json:"approved_by" db:"APPROVED_BY"`
	ApprovedStatus string    `json:"approved_status" db:"APPROVED_STAT"`
	AriaCityCode   int       `json:"aria_city_code" db:"ARIA_CITY_CODE"`
}
