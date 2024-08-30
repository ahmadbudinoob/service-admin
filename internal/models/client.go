package models

import (
	"time"
)

type UserClient struct {
	LoginID  string    `json:"login_id" db:"LOGIN_ID"`
	ClientCD string    `json:"client_cd" db:"CLIENT_CD"`
	CreateDT time.Time `json:"create_dt" db:"CREATE_DT"`
	CreateBy string    `json:"create_by" db:"CREATE_BY"`
}

type ClientDetail struct {
	ClientCD   string `json:"client_cd" db:"CLIENT_CD"`
	ClientName string `json:"client_nm" db:"CLIENT_NAME"`
}
