package main

import (
	"embed"
	"go-back/database"
	"log"
	"os"

	//"go-back/logic"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	database.InitDatabases()
	defer database.CloseDatabases()

	logEnabled, err := database.CredentialsDB.GetENVValue("logEnabled")

	var logFile *os.File

	if err == nil && logEnabled == "true" {
		logFile, err = os.OpenFile("output.log", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Ошибка открытия файла: %v", err)
		} else {
			log.SetOutput(logFile)
			defer logFile.Close()
		}
	}
	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Donate Agreagator",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app, //roll,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
