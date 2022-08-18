package gorush

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"git.raoulh.pw/raoulh/my_greenhouse_backend/config"
	logger "git.raoulh.pw/raoulh/my_greenhouse_backend/log"
	"github.com/sirupsen/logrus"
)

type GorushNotif struct {
	Notifications []interface{}
}

type AndroidNotif struct {
	Tokens       []string
	Platform     uint
	Message      string
	Title        string
	Notification struct {
		Icon  string
		Sound string
	}
}

type IOSNotif struct {
	Tokens      []string
	Platform    uint
	Message     string
	Topic       string
	Development bool
}

var logging *logrus.Entry

func init() {
	logging = logger.NewLogger("gorush")
}

func SendPushMessage(hwType uint, token string, message string, isDev bool) (err error) {
	url := config.Config.String("gorush.push_url")

	var data []byte
	if hwType == 1 {
		n := createIOSNotif(token, message, isDev)
		data, err = json.Marshal(n)
		logging.Errorf("json marshall failure: %s", err)
		if err != nil {
			return
		}
	} else if hwType == 2 {
		n := createAndroidNotif(token, message, isDev)
		data, err = json.Marshal(n)
		logging.Errorf("json marshall failure: %s", err)
		if err != nil {
			return
		}
	} else {
		return fmt.Errorf("unknown hw type")
	}

	body := bytes.NewBuffer(data)

	res, err := http.Post(url, "application/json", body)
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		defer res.Body.Close()
		b, _ := ioutil.ReadAll(res.Body)
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
		Topic:       "fr.calaos.myGreenhous",
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
