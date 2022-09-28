package models

import "time"

const (
	NotifHwIOS = iota + 1
	NotifHwAndroid
)

type User struct {
	ID                         uint                `gorm:"primarykey" json:"-"`
	CreatedAt                  time.Time           `json:"-"`
	UpdatedAt                  time.Time           `json:"-"`
	LastLogin                  time.Time           `json:"-"`
	DeviceID                   string              `json:"device_id" gorm:"index:idx_deviceid"`
	NotifToken                 string              `json:"-"`
	NotifHwType                uint                `json:"-"`
	NotifLocale                string              `json:"-"`
	NotifDevelopment           bool                `json:"-"`
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
	CurrentValue  float32   `json:"current_value"`
	CurrentTime   time.Time `json:"current_time"`
	HourAverage   float32   `json:"hour_average_value"`
	DayAverage    float32   `json:"day_average_value"`
	LastDayTime   time.Time `json:"last_day_time"`
	LastCheckTime time.Time `json:"-"` //Last time we did check for notification for this measurement
	LastValue     float32   `json:"-"` //Last value from api to compare from
}

type UnitMeasurements struct {
	ID                uint      `gorm:"primarykey" json:"-"`
	CreatedAt         time.Time `json:"-"`
	UpdatedAt         time.Time `json:"-"`
	UserID            uint      `json:"-"`
	ProdUnitID        uint      `json:"product_unit_id"`
	ProdUnitType      string    `json:"production_unit_type"`
	ProdUnitReference string    `json:"production_ref"`

	PH       Measurement `gorm:"embedded;embeddedPrefix:ph_" json:"ph"`
	Water    Measurement `gorm:"embedded;embeddedPrefix:water_" json:"watertemp"`
	Air      Measurement `gorm:"embedded;embeddedPrefix:air_" json:"airtemp"`
	Humidity Measurement `gorm:"embedded;embeddedPrefix:humidity_" json:"humidity"`
}

const (
	NotifTypePh uint = iota + 1
	NotifTypeWaterTemp
	NotifTypeAirTemp
	NotifTypeHumidity
)

type NotifSettings struct {
	ID         uint      `gorm:"primarykey" json:"-"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
	UserID     uint      `json:"-"`
	ProdUnitID uint      `json:"product_unit_id"`

	Type           uint          `json:"type"`
	RangeEnabled   bool          `json:"range_enabled"`
	RangeMin       float32       `json:"range_min"`
	RangeMax       float32       `json:"range_max"`
	TooFastEnabled bool          `json:"too_fast_enabled"`
	TimeEnabled    bool          `json:"time_enabled"`
	MinTime        time.Duration `json:"time_min"`
	DiffTime       time.Duration `json:"-" gorm:"-"`

	CurrentValue float32 `json:"-" gorm:"-"`
}
