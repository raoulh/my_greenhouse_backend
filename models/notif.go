package models

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"git.raoulh.pw/raoulh/my_greenhouse_backend/config"
	"git.raoulh.pw/raoulh/my_greenhouse_backend/models/gorush"
)

func GetNotifSettings(u *User, notifType uint) (n *NotifSettings, err error) {
	if u.ID == 0 {
		return n, fmt.Errorf("userID is 0")
	}

	n = &NotifSettings{}
	err = db.Where("user_id = ? and type = ?", u.ID, notifType).First(n).Error

	//this notif settings does not exist yet
	if err != nil {
		n = createDefaultNotifSettings(notifType)
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
	err = db.Where("user_id = ? and type = ?", u.ID, n.Type).First(notif).Error

	notif.UserID = u.ID
	notif.Type = n.Type
	notif.RangeEnabled = n.RangeEnabled
	notif.RangeMin = n.RangeMin
	notif.RangeMax = n.RangeMax
	notif.TooFastEnabled = n.TooFastEnabled
	notif.TimeEnabled = n.TimeEnabled
	notif.MinTime = n.MinTime

	if err != nil {
		err = db.Create(notif).Error
	} else {
		err = db.Save(notif).Error
	}

	return
}

func createDefaultNotifSettings(notifType uint) (n *NotifSettings) {
	n = &NotifSettings{
		Type:    notifType,
		MinTime: time.Hour * 2, //2 hours without data
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

	for idx := range u.Meas {
		n, err := GetNotifSettings(u, NotifTypePh)
		if err == nil {

			//Check out-of-range
			if n.RangeEnabled {
				if u.Meas[idx].PH.CurrentValue < n.RangeMin ||
					u.Meas[idx].PH.CurrentValue > n.RangeMax {
					n.CurrentValue = u.Meas[idx].PH.CurrentValue
					msg := u.executeTemplate(t, "ph_range", n)
					u.sendNotif(msg)
				}
			}

			if n.TooFastEnabled {
				//TODO
			}

			if n.TimeEnabled {
				//TODO
			}

		}

		n, err = GetNotifSettings(u, NotifTypeWaterTemp)
		if err == nil {

			//Check out-of-range
			if n.RangeEnabled {
				if u.Meas[idx].Water.CurrentValue < n.RangeMin ||
					u.Meas[idx].Water.CurrentValue > n.RangeMax {
					n.CurrentValue = u.Meas[idx].Water.CurrentValue
					msg := u.executeTemplate(t, "water_range", n)
					u.sendNotif(msg)
				}
			}

			if n.TooFastEnabled {
				//TODO
			}

			if n.TimeEnabled {
				//TODO
			}

		}

		n, err = GetNotifSettings(u, NotifTypeAirTemp)
		if err == nil {

			//Check out-of-range
			if n.RangeEnabled {
				if u.Meas[idx].Air.CurrentValue < n.RangeMin ||
					u.Meas[idx].Air.CurrentValue > n.RangeMax {
					n.CurrentValue = u.Meas[idx].Air.CurrentValue
					msg := u.executeTemplate(t, "air_range", n)
					u.sendNotif(msg)
				}
			}

			if n.TooFastEnabled {
				//TODO
			}

			if n.TimeEnabled {
				//TODO
			}

		}

		n, err = GetNotifSettings(u, NotifTypeHumidity)
		if err == nil {

			//Check out-of-range
			if n.RangeEnabled {
				if u.Meas[idx].Humidity.CurrentValue < n.RangeMin ||
					u.Meas[idx].Humidity.CurrentValue > n.RangeMax {
					n.CurrentValue = u.Meas[idx].Humidity.CurrentValue
					msg := u.executeTemplate(t, "humidity_range", n)
					u.sendNotif(msg)
				}
			}

			if n.TooFastEnabled {
				//TODO
			}

			if n.TimeEnabled {
				//TODO
			}

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
