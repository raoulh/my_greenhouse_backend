package gorush

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"git.raoulh.pw/raoulh/my_greenhouse_backend/config"
	logger "git.raoulh.pw/raoulh/my_greenhouse_backend/log"
	"github.com/sirupsen/logrus"
)

type GorushNotif struct {
	Notifications []interface{}
}

type AndroidNotif struct {
	Tokens       []string `json:"tokens"`
	Platform     uint     `json:"platform"`
	Message      string   `json:"message"`
	Title        string   `json:"title"`
	Notification struct {
		Icon  string `json:"icon"`
		Sound string `json:"sound"`
	} `json:"notification"`
}

type IOSNotif struct {
	Tokens      []string `json:"tokens"`
	Platform    uint     `json:"platform"`
	Message     string   `json:"message"`
	Topic       string   `json:"topic"`
	Development bool     `json:"development"`
}

var logging *logrus.Entry

func init() {
	logging = logger.NewLogger("gorush")
}

func SendPushMessage(hwType uint, token string, message string, isDev bool) (err error) {
	url := config.Config.String("gorush.push_url")

	logging.Debugf("Sending push to gorush@: %s", url)

	gn := GorushNotif{}

	if hwType == 1 {
		n := createIOSNotif(token, message, isDev)
		gn.Notifications = append(gn.Notifications, n)
	} else if hwType == 2 {
		n := createAndroidNotif(token, message, isDev)
		gn.Notifications = append(gn.Notifications, n)
	} else {
		return fmt.Errorf("unknown hw type")
	}

	var data []byte
	data, err = json.Marshal(gn)
	if err != nil {
		logging.Errorf("json marshall failure: %s", err)
		return
	}

	logging.Debugf("JSON notif: %s", string(data))

	body := bytes.NewBuffer(data)

	res, err := http.Post(url, "application/json", body)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		defer res.Body.Close()
		b, _ := io.ReadAll(res.Body)
		logging.Errorf("gorush notification failed: %s", string(b))
		return fmt.Errorf("gorush failure")
	}

	return
}

func createIOSNotif(token string, message string, isDev bool) *IOSNotif {
	return &IOSNotif{
		Tokens:      []string{token},
		Platform:    1,
		Message:     message,
		Topic:       "fr.calaos.myGreenhouse",
		Development: isDev,
	}
}

func createAndroidNotif(token string, message string, isDev bool) *AndroidNotif {
	n := &AndroidNotif{
		Tokens:   []string{token},
		Platform: 2,
		Message:  message,
		Title:    "MyFood",
	}

	n.Notification.Icon = "ic_stat_myfood"
	n.Notification.Sound = "default"

	return n
}
