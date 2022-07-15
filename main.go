package main

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"git.raoulh.pw/raoulh/my_greenhouse_backend/app"
	"git.raoulh.pw/raoulh/my_greenhouse_backend/config"
	"git.raoulh.pw/raoulh/my_greenhouse_backend/models"

	logger "git.raoulh.pw/raoulh/my_greenhouse_backend/log"

	"github.com/fatih/color"
	cli "github.com/jawher/mow.cli"
	"github.com/sirupsen/logrus"
)

const (
	DefaultConfigFilename = "/etc/greenhouse.toml"

	CharStar     = "\u2737"
	CharAbort    = "\u2718"
	CharCheck    = "\u2714"
	CharWarning  = "\u26A0"
	CharArrow    = "\u2012\u25b6"
	CharVertLine = "\u2502"
)

var (
	blue       = color.New(color.FgBlue).SprintFunc()
	errorRed   = color.New(color.FgRed).SprintFunc()
	errorBgRed = color.New(color.BgRed, color.FgBlack).SprintFunc()
	green      = color.New(color.FgGreen).SprintFunc()
	cyan       = color.New(color.FgCyan).SprintFunc()
	bgCyan     = color.New(color.FgWhite).SprintFunc()

	logging *logrus.Entry

	myApp *app.AppServer
)

func exit(err error, exit int) {
	logging.Fatalln(errorRed(CharAbort), err)
	cli.Exit(exit)
}

func handleSignals() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sigint

	logging.Println("Shuting down...")
	myApp.Shutdown()
	models.Shutdown()
}

func main() {
	logging = logger.NewLogger("greehouse")
	runtime.GOMAXPROCS(runtime.NumCPU())

	a := cli.App("greehouse_backend", "Greenhouse Backend")

	a.Spec = "[-c]"

	var (
		conffile = a.StringOpt("c config", DefaultConfigFilename, "Set config file")
	)

	a.Action = func() {
		var err error

		if err = config.InitConfig(conffile); err != nil {
			exit(err, 1)
		}

		if myApp, err = app.NewApp(); err != nil {
			exit(err, 1)
		}

		if err = models.Init(); err != nil {
			exit(err, 1)
		}

		myApp.Start()

		handleSignals()
	}

	if err := a.Run(os.Args); err != nil {
		exit(err, 1)
	}
}
