package models

import (
	"fmt"
	"time"
)

func GetNotifSettings(u *User, notifType uint) (n *NotifSettings, err error) {
	n = &NotifSettings{}
	err = db.Where("user_id = ? and type = ?", u.ID, notifType).First(n).Error

	//this notif settings does not exist yet
	if err == nil {
		n = createDefaultNotifSettings(notifType)
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

	notif.Type = n.Type
	notif.RangeEnabled = n.RangeEnabled
	notif.RangeMin = n.RangeMin
	notif.RangeMax = n.RangeMax
	notif.TooFastEnabled = n.TooFastEnabled
	notif.TimeEnabled = n.TimeEnabled
	notif.MinTime = n.MinTime

	if err == nil {
		err = db.Create(n).Error
	} else {
		err = db.Save(n).Error
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
