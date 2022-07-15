package myfood

import (
	"encoding/json"
	"time"
)

type BaseMessage struct {
	Failed    bool     `json:"failed"`
	Messages  []string `json:"messages"`
	Succeeded bool     `json:"succeeded"`
}

type AuthData struct {
	Username string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
}

type TokenData struct {
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type TokenTime struct {
	time.Time
}

// UnmarshalJSON decodes 0001-01-01T00:00:00 into a time.Time object
func (p *TokenTime) UnmarshalJSON(bytes []byte) error {
	var raw string
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02T15:04:05", raw)
	if err != nil {
		return err
	}

	p.Time = t
	return nil
}

type TokenResultData struct {
	BaseMessage
	Data struct {
		TokenData
		RefreshTokenExpiryTime TokenTime `json:"refreshTokenExpiryTime,omitempty"`
	}
}

type LastDayDate struct {
	time.Time
}

// UnmarshalJSON decodes 0001-01-01T00:00:00 into a time.Time object
func (p *LastDayDate) UnmarshalJSON(bytes []byte) error {
	var raw string
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	t, err := time.Parse("1/2/2006", raw)
	if err != nil {
		return err
	}

	p.Time = t
	return nil
}

type LastCaptureTime struct {
	time.Time
}

// UnmarshalJSON decodes 0001-01-01T00:00:00 into a time.Time object
func (p *LastCaptureTime) UnmarshalJSON(bytes []byte) error {
	var raw string
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	t, err := time.Parse("3:04 PM", raw)
	if err != nil {
		return err
	}

	p.Time = t
	return nil
}

type ProdUnitDetailData struct {
	BaseMessage
	Data struct {
		PioneerCitizenName          string          `json:"pioneerCitizenName"`
		PioneerCitizenNumber        uint            `json:"pioneerCitizenNumber"`
		ProductionUnitVersion       string          `json:"productionUnitVersion"`
		ProductionUnitType          string          `json:"productionUnitType"`
		PicturePath                 string          `json:"picturePath"`
		ProductionUnitOptions       string          `json:"productionUnitOptions"`
		OnlineSinceWeeks            uint            `json:"onlineSinceWeeks"`
		AverageMonthlyProduction    uint            `json:"averageMonthlyProduction"`
		AverageMonthlySparedCO2     float32         `json:"averageMonthlySparedCO2"`
		CurrentPhValue              float32         `json:"currentPhValue"`
		CurrentPhCaptureTime        LastCaptureTime `json:"currentPhCaptureTime"`
		AverageHourPhValue          float32         `json:"averageHourPhValue"`
		AverageDayPhValue           float32         `json:"averageDayPhValue"`
		LastDayPhCaptureTime        LastDayDate     `json:"lastDayPhCaptureTime"`
		CurrentWaterTempValue       float32         `json:"currentWaterTempValue"`
		CurrentWaterTempCaptureTime LastCaptureTime `json:"currentWaterTempCaptureTime"`
		AverageHourWaterTempValue   float32         `json:"averageHourWaterTempValue"`
		AverageDayWaterTempValue    float32         `json:"averageDayWaterTempValue"`
		LastDayWaterTempCaptureTime LastDayDate     `json:"lastDayWaterTempCaptureTime"`
		CurrentAirTempValue         float32         `json:"currentAirTempValue"`
		CurrentAirTempCaptureTime   LastCaptureTime `json:"currentAirTempCaptureTime"`
		AverageHourAirTempValue     float32         `json:"averageHourAirTempValue"`
		AverageDayAirTempValue      float32         `json:"averageDayAirTempValue"`
		LastDayAirTempCaptureTime   LastDayDate     `json:"lastDayAirTempCaptureTime"`
		CurrentHumidityValue        float32         `json:"currentHumidityValue"`
		CurrentHumidityCaptureTime  LastCaptureTime `json:"currentHumidityCaptureTime"`
		AverageHourHumidityValue    float32         `json:"averageHourHumidityValue"`
		AverageDayHumidityValue     float32         `json:"averageDayHumidityValue"`
		LastDayHumidityCaptureTime  LastDayDate     `json:"lastDayHumidityCaptureTime"`
		LastSignalStrenghtReceived  string          `json:"lastSignalStrenghtReceived"`
	}
}

type ProdUnitsData struct {
	BaseMessage
	Data struct {
		ProductUnits []struct {
			ID        uint   `json:"id"`
			Reference string `json:"reference"`
			Info      string `json:"info"`
		} `json:"productUnits"`
	}
}

type ResultDataTime struct {
	time.Time
}

// UnmarshalJSON decodes 0001-01-01T00:00:00 into a time.Time object
func (p *ResultDataTime) UnmarshalJSON(bytes []byte) error {
	var raw string
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02T15:04:05.999999999", raw)
	if err != nil {
		return err
	}

	p.Time = t
	return nil
}

type ResultData struct {
	BaseMessage
	Data struct {
		ResultData []struct {
			Value       float32        `json:"value"`
			CaptureDate ResultDataTime `json:"captureDate"`
		} `json:"resultData"`
	}
}
