package myfood

import (
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

type TokenResultData struct {
	BaseMessage
	Data struct {
		TokenData
		RefreshTokenExpiryTime time.Time `json:"refreshTokenExpiryTime,omitempty"`
	}
}

type ProdUnitDetailData struct {
	BaseMessage
	Data struct {
		PioneerCitizenName          string    `json:"pioneerCitizenName"`
		PioneerCitizenNumber        uint      `json:"pioneerCitizenNumber"`
		ProductionUnitVersion       string    `json:"productionUnitVersion"`
		ProductionUnitType          string    `json:"productionUnitType"`
		PicturePath                 string    `json:"picturePath"`
		ProductionUnitOptions       string    `json:"productionUnitOptions"`
		OnlineSinceWeeks            uint      `json:"onlineSinceWeeks"`
		AverageMonthlyProduction    string    `json:"averageMonthlyProduction"`
		AverageMonthlySparedCO2     float32   `json:"averageMonthlySparedCO2"`
		CurrentPhValue              float32   `json:"currentPhValue"`
		CurrentPhCaptureTime        time.Time `json:"currentPhCaptureTime"`
		AverageHourPhValue          float32   `json:"averageHourPhValue"`
		AverageDayPhValue           float32   `json:"averageDayPhValue"`
		LastDayPhCaptureTime        time.Time `json:"lastDayPhCaptureTime"`
		CurrentWaterTempValue       float32   `json:"currentWaterTempValue"`
		CurrentWaterTempCaptureTime time.Time `json:"currentWaterTempCaptureTime"`
		AverageHourWaterTempValue   float32   `json:"averageHourWaterTempValue"`
		AverageDayWaterTempValue    float32   `json:"averageDayWaterTempValue"`
		LastDayWaterTempCaptureTime time.Time `json:"lastDayWaterTempCaptureTime"`
		CurrentAirTempValue         float32   `json:"currentAirTempValue"`
		CurrentAirTempCaptureTime   time.Time `json:"currentAirTempCaptureTime"`
		AverageHourAirTempValue     float32   `json:"averageHourAirTempValue"`
		AverageDayAirTempValue      float32   `json:"averageDayAirTempValue"`
		LastDayAirTempCaptureTime   time.Time `json:"lastDayAirTempCaptureTime"`
		CurrentHumidityValue        float32   `json:"currentHumidityValue"`
		CurrentHumidityCaptureTime  time.Time `json:"currentHumidityCaptureTime"`
		AverageHourHumidityValue    float32   `json:"averageHourHumidityValue"`
		AverageDayHumidityValue     float32   `json:"averageDayHumidityValue"`
		LastDayHumidityCaptureTime  time.Time `json:"lastDayHumidityCaptureTime"`
		LastSignalStrenghtReceived  string    `json:"lastSignalStrenghtReceived"`
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
			Value       float32   `json:"value"`
			CaptureDate time.Time `json:"captureDate"`
		} `json:"resultData"`
	}
}
