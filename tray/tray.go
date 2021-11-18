package tray

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"os"
	"os/signal"
	"syscall"

	"github.com/getlantern/systray"
	"github.com/hirokimoto/crypto-auto/views"
	"github.com/skratchdot/open-golang/open"
)

func OnReady() {
	systray.SetIcon(getIcon("assets/auto.ico"))

	mHelloWorld := systray.AddMenuItem("Hello, World!", "Opens a simple HTML Hello, World")
	systray.AddSeparator()
	mGoogleBrowser := systray.AddMenuItem("Google in Browser", "Opens Google in a normal browser")
	mGoogleEmbed := systray.AddMenuItem("Google in Window", "Opens Google in a custom window")
	mSettings := systray.AddMenuItem("Settings in Window", "Opens Google in a custom window")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit example tray application")

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			systray.SetTitle(getClockTime("Local"))
			systray.SetTooltip("Local timezone")
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		select {

		case <-mHelloWorld.ClickedCh:
			err := views.Get().OpenIndex()
			if err != nil {
				fmt.Println(err)
			}
		case <-mGoogleBrowser.ClickedCh:
			err := open.Run("https://www.google.com")
			if err != nil {
				fmt.Println(err)
			}
		case <-mGoogleEmbed.ClickedCh:
			err := views.Get().OpenGoogle()
			if err != nil {
				fmt.Println(err)
			}
		case <-mSettings.ClickedCh:
			err := views.Get().OpenSettings()
			if err != nil {
				fmt.Println(err)
			}
		case <-mQuit.ClickedCh:
			systray.Quit()
		case <-sigc:
			systray.Quit()
		}
	}
}

func OnQuit() {
	close(views.Get().Shutdown)
}

func getClockTime(tz string) string {
	t := time.Now()
	utc, _ := time.LoadLocation(tz)

	hour, min, sec := t.In(utc).Clock()
	return itoaTwoDigits(hour) + ":" + itoaTwoDigits(min) + ":" + itoaTwoDigits(sec)
}

func itoaTwoDigits(i int) string {
	b := "0" + strconv.Itoa(i)
	return b[len(b)-2:]
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}
