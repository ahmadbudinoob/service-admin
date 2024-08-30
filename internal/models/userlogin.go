package models

import (
	"time"
)

// TLDUserLoginLog represents the TLD_USERLOGIN_LOG table
type UserLogin struct {
	LoginID       string    `json:"login_id" db:"LOGINID"`
	Status        string    `json:"status" db:"STATUS"`
	ActionDate    time.Time `json:"action_date" db:"ACTION_DATE"`
	ActionTime    time.Time `json:"action_time" db:"ACTION_TIME"`
	Channel       string    `json:"channel,omitempty" db:"CHANNEL"`
	ChannelMedia  string    `json:"channel_media,omitempty" db:"CHANNEL_MEDIA"`
	ChannelDevice string    `json:"channel_device,omitempty" db:"CHANNEL_DEVICE"`
	ScreenWidth   int       `json:"screen_width,omitempty" db:"SCREEN_WIDTH"`
	ScreenHeight  int       `json:"screen_height,omitempty" db:"SCREEN_HEIGHT"`
	IPAddress     string    `json:"ip_address,omitempty" db:"IP_ADDRESS"`
}
