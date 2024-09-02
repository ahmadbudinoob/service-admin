package models

import (
	"database/sql"
	"time"
)

// User struct maps to the TLOTSUSER table
type User struct {
	LoginID           string         `db:"LOGIN_ID" json:"login_id"`
	FullName          string         `db:"FULL_NAME" json:"full_name"`
	Password          string         `db:"PASSWD" json:"password"`
	PasswordExpDate   time.Time      `db:"PASSWD_EXPDATE" json:"password_exp_date"`
	UserStatus        string         `db:"USER_STATUS" json:"user_status"`
	DescStatus        string         `db:"DESC_STATUS" json:"desc_status,omitempty"`
	OrderRestrictions string         `db:"ORDERRESTRICTIONS" json:"order_restrictions"`
	PIN               string         `db:"PIN" json:"pin"`
	PINExpDate        time.Time      `db:"PIN_EXPDATE" json:"pin_exp_date"`
	MasterClientCD    sql.NullString `db:"MASTER_CLIENT_CD" json:"master_client_cd,omitempty"`
	LastLogin         sql.NullTime   `db:"LAST_LOGIN" json:"last_login,omitempty"`
	CreateBy          string         `db:"CREATE_BY" json:"create_by"`
	CreateDT          time.Time      `db:"CREATE_DT" json:"create_dt"`
	UpdateBy          string         `db:"UPDATE_BY" json:"update_by"`
	UpdateDT          time.Time      `db:"UPDATE_DT" json:"update_dt"`
	PhotoID           int32          `db:"PHOTO_ID" json:"photo_id,omitempty"`
	IsProtlAm         sql.NullString `db:"IS_PROTLAM" json:"is_protl_am,omitempty"`
	Email             sql.NullString `db:"EMAIL" json:"email,omitempty"`
	ClientBirthDT     sql.NullTime   `db:"CLIENT_BIRTH_DT" json:"client_birth_dt,omitempty"`
	City              sql.NullInt16  `db:"CITY" json:"city,omitempty"`
	Telepon           sql.NullString `db:"TELEPON" json:"telepon,omitempty"`
}

type CreateUserRequest struct {
	LoginID           string     `json:"login_id" validate:"required"`
	FullName          string     `json:"full_name" validate:"required"`
	Password          string     `json:"password" validate:"required"`
	OrderRestrictions string     `json:"order_restrictions" validate:"required"`
	Pin               string     `json:"pin" validate:"required"`
	ExpiresAt         *time.Time `json:"expires_at,omitempty"`
}

type UserResponseID struct {
	LoginID           string     `json:"login_id"`
	FullName          string     `json:"full_name"`
	UserStatus        string     `json:"user_status"`
	OrderRestrictions string     `json:"order_restrictions"`
	Password          *string    `json:"password,omitempty"`
	PasswordExpDate   *time.Time `json:"password_exp_date,omitempty"`
	Pin               *string    `json:"pin,omitempty"`
	PinExpDate        *time.Time `json:"pin_exp_date,omitempty"`
	MasterClientCD    *string    `json:"master_client_cd,omitempty"`
	LastLogin         *time.Time `json:"last_login,omitempty"`
	CreateBy          string     `json:"create_by"`
	CreateDT          time.Time  `json:"create_dt"`
	UpdateBy          string     `json:"update_by"`
	UpdateDT          time.Time  `json:"update_dt"`
	PhotoID           *int16     `json:"photo_id,omitempty"`
	IsProtlAm         *string    `json:"is_protl_am,omitempty"`
	Email             *string    `json:"email,omitempty"`
	ClientBirthDT     *time.Time `json:"client_birth_dt,omitempty"`
	City              *int16     `json:"city,omitempty"`
	Telepon           *string    `json:"telepon,omitempty"`
}

type UserResponse struct {
	LoginID           string `json:"LoginID"`
	FullName          string `json:"FullName"`
	UserStatus        string `json:"UserStatus"`
	OrderRestrictions string `json:"OrderRestrictions"`
	CreateDT          string `json:"CreateDT"`
	UpdateDT          string `json:"UpdateDT"`
}

type UpdateUserRequest struct {
	LoginID           string `json:"login_id"`                          // Mandatory
	FullName          string `json:"full_name"`                         // Optional
	Phone             string `json:"phone"`                             // Optional
	Email             string `json:"email"`                             // Optional
	ClientBirthDate   string `json:"client_birth_dt" format:"dateTime"` // Optional
	City              int    `json:"city"`                              // Optional
	OrderRestrictions string `json:"order_restrictions"`                // Optional
	UserStatus        string `json:"user_status"`
}
