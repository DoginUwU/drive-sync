package main

import (
	"os"

	"github.com/getlantern/systray"
)

func startTray() {
	systray.Run(onTrayReady, onTrayExit)
}

func onTrayReady() {
	systray.SetIcon(getIcon("assets/icon.ico"))
	systray.SetTitle("Google Drive Sync")
	systray.SetTooltip("Last sync in: 10:34")

	addTrayItems()
}

func addTrayItems() {
	systray.AddMenuItem("doginuwu@gmail.com", "Last sync in: 10:34")
	systray.AddSeparator()

	settingsButton := systray.AddMenuItem("Open settings", "")

	systray.AddSeparator()
	quitButton := systray.AddMenuItem("Quit", "Close the Google Drive Sync")

	go func() {
		for {
			select {
			case <-settingsButton.ClickedCh:
				showSettings()
			case <-quitButton.ClickedCh:
				systray.Quit()
				quitSettings()
				os.Exit(0)
			}
		}
	}()
}

func onTrayExit() {

}
