package models

import (
	"bytes"
	"fmt"
	"math"
	"text/template"
	"time"

	"git.raoulh.pw/raoulh/my_greenhouse_backend/config"
	"git.raoulh.pw/raoulh/my_greenhouse_backend/models/gorush"
)

func GetNotifSettings(u *User, notifType, prodId uint) (n *NotifSettings, err error) {
	if u.ID == 0 {
		return n, fmt.Errorf("userID is 0")
	}

	n = &NotifSettings{}
	err = db.Where("user_id = ? and type = ? and prod_unit_id = ?", u.ID, notifType, prodId).First(n).Error

	//this notif settings does not exist yet
	if err != nil {
		n = createDefaultNotifSettings(notifType, prodId)
		n.UserID = u.ID
		err = db.Create(n).Error
		return
	}

	return
}

func SetNotifSettings(u *User, n *NotifSettings) (err error) {
	if n.Type == 0 {
		return fmt.Errorf("bad NotifSettings")
	}

	notif := &NotifSettings{}
	err = db.Where("user_id = ? and type = ? and prod_unit_id = ?", u.ID, n.Type, n.ProdUnitID).First(notif).Error

	notif.UserID = u.ID
	notif.Type = n.Type
	notif.RangeEnabled = n.RangeEnabled
	notif.RangeMin = n.RangeMin
	notif.RangeMax = n.RangeMax
	notif.TooFastEnabled = n.TooFastEnabled
	notif.TimeEnabled = n.TimeEnabled
	notif.MinTime = n.MinTime

	if n.ProdUnitID <= 0 {
		return fmt.Errorf("bad product_unit_id")
	}

	if err != nil {
		err = db.Create(notif).Error
	} else {
		err = db.Save(notif).Error
	}

	return
}

func createDefaultNotifSettings(notifType, prodId uint) (n *NotifSettings) {
	n = &NotifSettings{
		Type:       notifType,
		MinTime:    time.Hour * 2, //2 hours without data
		ProdUnitID: prodId,
	}

	switch notifType {
	case NotifTypePh:
		n.RangeMin = 5
		n.RangeMax = 8
	case NotifTypeAirTemp:
		n.RangeMin = 8
		n.RangeMax = 38
	case NotifTypeWaterTemp:
		n.RangeMin = 10
		n.RangeMax = 30
	case NotifTypeHumidity:
		n.RangeMin = 15
		n.RangeMax = 80
	}

	return
}

func UpdateNotifToken(u *User, token string, hwType uint, locale string, dev bool) (err error) {
	u.NotifToken = token
	u.NotifHwType = hwType
	u.NotifLocale = locale
	u.NotifDevelopment = dev

	err = db.Save(u).Error
	return
}

// Check if any notifications should be sent to user based on his configured options
func (u *User) handleNotifications() {
	t := template.Must(template.ParseGlob(config.Config.String("general.data") + "/*.notif"))

	for idx, meas := range u.Meas {
		n, err := GetNotifSettings(u, NotifTypePh, meas.ProdUnitID)
		if err == nil {
			u.doCheckNotif(&u.Meas[idx].PH, n, t, "ph_range", "ph_fast", "ph_time", 0.8)
		}

		n, err = GetNotifSettings(u, NotifTypeWaterTemp, meas.ProdUnitID)
		if err == nil {
			u.doCheckNotif(&u.Meas[idx].Water, n, t, "water_range", "water_fast", "water_time", 10)
		}

		n, err = GetNotifSettings(u, NotifTypeAirTemp, meas.ProdUnitID)
		if err == nil {
			u.doCheckNotif(&u.Meas[idx].Air, n, t, "air_range", "air_fast", "air_time", 10)
		}

		n, err = GetNotifSettings(u, NotifTypeHumidity, meas.ProdUnitID)
		if err == nil {
			u.doCheckNotif(&u.Meas[idx].Humidity, n, t, "humidity_range", "humidity_fast", "humidity_time", 30)
		}
	}
}

func (u *User) doCheckNotif(meas *Measurement, n *NotifSettings, t *template.Template, rangeTpl, toofastTpl, timeTpl string, threshold float32) {
	hasNewValue := !meas.LastCheckTime.Equal(meas.CurrentTime)
	if hasNewValue {
		meas.LastCheckTime = meas.CurrentTime
	}

	//Check out-of-range
	if n.RangeEnabled &&
		hasNewValue {
		if meas.CurrentValue < n.RangeMin ||
			meas.CurrentValue > n.RangeMax {
			n.CurrentValue = meas.CurrentValue
			msg := u.executeTemplate(t, rangeTpl, n)
			u.sendNotif(msg)
		}
	}

	//Check if value changes too fast
	if n.TooFastEnabled &&
		hasNewValue {
		if meas.LastValue > 0 {
			diff := meas.CurrentValue - meas.LastValue
			if float32(math.Abs(float64(diff))) < threshold {
				msg := u.executeTemplate(t, toofastTpl, n)
				u.sendNotif(msg)
			}
		}
		meas.LastValue = meas.CurrentValue
	}

	//Check if data has not been received since duration
	if n.TimeEnabled {
		hoursDiff := time.Since(meas.LastCheckTime)
		if hoursDiff.Hours() > n.MinTime.Hours() {
			n.DiffTime = hoursDiff
			msg := u.executeTemplate(t, timeTpl, n)
			u.sendNotif(msg)
		}
	}
}

func (u *User) executeTemplate(t *template.Template, templateName string, n *NotifSettings) string {
	var b bytes.Buffer

	tt := t.Lookup(fmt.Sprintf("%s.%s.notif", templateName, u.NotifLocale))
	if tt == nil {
		tt = t.Lookup(fmt.Sprintf("%s.en_US.notif", templateName))

		if tt == nil {
			logging.Errorf("template named '%s' does not exist", templateName)
			return ""
		}
	}

	if err := tt.Execute(&b, n); err != nil {
		logging.Errorf("Failed to execute notification template: %v", err)
		return ""
	}

	return b.String()
}

func (u *User) sendNotif(msg string) {
	if msg == "" {
		return
	}
	if err := gorush.SendPushMessage(u.NotifHwType, u.NotifToken, msg, u.NotifDevelopment); err != nil {
		logging.Errorf("Failed to send gorush notif: %v", err)
	}
}
