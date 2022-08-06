package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var application fyne.App
var window fyne.Window

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
	loginBox := container.NewHBox(
		canvas.NewText("Account Login:", color.White),
		widget.NewButton("Synced with doginuwu@gmail.com", syncAccount),
	)

	container := container.NewVBox()

	container.Add(loginBox)

	return container
}

func syncAccount() {

}

func showSettings() {
	window.Show()
}

func quitSettings() {
	application.Quit()
}
