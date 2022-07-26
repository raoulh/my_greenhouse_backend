package models

import "time"

type User struct {
	ID                         uint                `gorm:"primarykey" json:"-"`
	CreatedAt                  time.Time           `json:"-"`
	UpdatedAt                  time.Time           `json:"-"`
	LastLogin                  time.Time           `json:"-"`
	DeviceID                   string              `json:"device_id" gorm:"index:idx_deviceid"`
	MF_Username                string              `json:"myfood_username"`
	MF_Token                   string              `json:"myfood_token,omitempty" gorm:"index:idx_token,unique"`
	MF_RefreshToken            string              `json:"-"`
	MF_TokenValidity           time.Time           `json:"myfood_validity"`
	MF_TokenValid              bool                `json:"myfood_token_valid"`
	MF_UserID                  uint                `json:"pioneer_citizen_number"`
	MF_RefreshTokenFailureCout uint                `json:"-" gorm:"refresh_token_count"`
	Meas                       []*UnitMeasurements `json:"meas"`
}

type Measurement struct {
	CurrentValue float32   `json:"current_value"`
	CurrentTime  time.Time `json:"current_time"`
	HourAverage  float32   `json:"hour_average_value"`
	DayAverage   float32   `json:"day_average_value"`
	LastDayTime  time.Time `json:"last_day_time"`
}

type UnitMeasurements struct {
	ID           uint      `gorm:"primarykey" json:"-"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
	UserID       uint      `json:"-"`
	ProdUnitID   uint      `json:"product_unit_id"`
	ProdUnitType string    `json:"production_unit_type"`

	PH       Measurement `gorm:"embedded;embeddedPrefix:ph_" json:"ph"`
	Water    Measurement `gorm:"embedded;embeddedPrefix:water_" json:"watertemp"`
	Air      Measurement `gorm:"embedded;embeddedPrefix:air_" json:"airtemp"`
	Humidity Measurement `gorm:"embedded;embeddedPrefix:humidity_" json:"humidity"`
}
