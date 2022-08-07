package main

import (
	"fmt"
	"os"

	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var application fyne.App
var window fyne.Window
var directoryButton *widget.Button
var syncButton *widget.Button

func startSettings() {
	application = app.New()
	window = application.NewWindow("Settings")
	window.Resize(fyne.NewSize(800, 400))

	createSettingsContent()

	application.Run()
}

func createSettingsContent() {
	window.SetContent(createSettingsLayout())
}

func createSettingsLayout() fyne.CanvasObject {
	syncButton = widget.NewButton("Synced with jhon_doe@gmail.com", syncAccount)

	loginBox := container.NewHBox(
		canvas.NewText("Account Login:", color.White),
		syncButton,
	)

	directoryButton = widget.NewButton(getDirectory(), searchDirectory)

	directoryBox := container.NewHBox(
		canvas.NewText("Google Drive directory:", color.White),
		directoryButton,
	)

	syncNowButton := widget.NewButton("Sync Now", startSync)

	container := container.NewVBox()

	container.Add(loginBox)
	container.Add(directoryBox)
	container.Add(syncNowButton)

	return container
}

func syncAccount() {
	fmt.Println("Syncing account...")
	os.Remove("token.json")
	startAuth()
}

func searchDirectory() {
	dialog.ShowFolderOpen(
		func(dir fyne.ListableURI, err error) {
			if err == nil {
				setDirectory(dir.Path())
			}

			directoryButton.SetText(getDirectory())
		},
		window,
	)
}

func showSettings() {
	window.Show()
}

func quitSettings() {
	application.Quit()
}
