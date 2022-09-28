package myfood

import (
	"encoding/json"
	"fmt"
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

type ApiTime struct {
	time.Time
}

type TokenResultData struct {
	BaseMessage
	Data struct {
		TokenData
		RefreshTokenExpiryTime ApiTime `json:"refreshTokenExpiryTime,omitempty"`
	}
}

// UnmarshalJSON decodes all times format from API into a time.Time object
func (p *ApiTime) UnmarshalJSON(bytes []byte) error {
	var raw string
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		return err
	}

	if raw == "-" {
		return nil
	}

	formats := []string{"3:04 PM", "1/2/2006", "2006-01-02T15:04:05Z", "2006-01-02T15:04:05-07:00", "2006-01-02T15:04:05", "2006-01-02T15:04:05.999999999"}

	parseOk := false
	for _, f := range formats {
		t, err := time.Parse(f, raw)
		if err == nil {
			p.Time = t
			parseOk = true
		}
	}

	if !parseOk {
		return fmt.Errorf("failed to parse time for %s", raw)
	}

	return nil
}

type ProdUnitDetailData struct {
	BaseMessage
	Data struct {
		PioneerCitizenName          string  `json:"pioneerCitizenName"`
		PioneerCitizenNumber        uint    `json:"pioneerCitizenNumber"`
		ProductionUnitVersion       string  `json:"productionUnitVersion"`
		ProductionUnitType          string  `json:"productionUnitType"`
		PicturePath                 string  `json:"picturePath"`
		ProductionUnitOptions       string  `json:"productionUnitOptions"`
		OnlineSinceWeeks            uint    `json:"onlineSinceWeeks"`
		AverageMonthlyProduction    uint    `json:"averageMonthlyProduction"`
		AverageMonthlySparedCO2     float32 `json:"averageMonthlySparedCO2"`
		CurrentPhValue              float32 `json:"currentPhValue"`
		CurrentPhCaptureTime        ApiTime `json:"currentPhCaptureTime"`
		AverageHourPhValue          float32 `json:"averageHourPhValue"`
		AverageDayPhValue           float32 `json:"averageDayPhValue"`
		LastDayPhCaptureTime        ApiTime `json:"lastDayPhCaptureTime"`
		CurrentWaterTempValue       float32 `json:"currentWaterTempValue"`
		CurrentWaterTempCaptureTime ApiTime `json:"currentWaterTempCaptureTime"`
		AverageHourWaterTempValue   float32 `json:"averageHourWaterTempValue"`
		AverageDayWaterTempValue    float32 `json:"averageDayWaterTempValue"`
		LastDayWaterTempCaptureTime ApiTime `json:"lastDayWaterTempCaptureTime"`
		CurrentAirTempValue         float32 `json:"currentAirTempValue"`
		CurrentAirTempCaptureTime   ApiTime `json:"currentAirTempCaptureTime"`
		AverageHourAirTempValue     float32 `json:"averageHourAirTempValue"`
		AverageDayAirTempValue      float32 `json:"averageDayAirTempValue"`
		LastDayAirTempCaptureTime   ApiTime `json:"lastDayAirTempCaptureTime"`
		CurrentHumidityValue        float32 `json:"currentHumidityValue"`
		CurrentHumidityCaptureTime  ApiTime `json:"currentHumidityCaptureTime"`
		AverageHourHumidityValue    float32 `json:"averageHourHumidityValue"`
		AverageDayHumidityValue     float32 `json:"averageDayHumidityValue"`
		LastDayHumidityCaptureTime  ApiTime `json:"lastDayHumidityCaptureTime"`
		LastSignalStrenghtReceived  string  `json:"lastSignalStrenghtReceived"`
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

type ResultData struct {
	BaseMessage
	Data struct {
		ResultData []struct {
			Value       float32 `json:"value"`
			CaptureDate ApiTime `json:"captureDate"`
		} `json:"resultData"`
	}
}
