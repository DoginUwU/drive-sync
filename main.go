package main

func main() {
	go startTray()
	go startBackend()
	go startAuth()
	startSettings()
}
