package main

type global struct {
	Text   TextLib
	Config Config
	Font   Font
	G      gameData
}

var globals global

func initGlobalSetting() {
	globals.Text = loadText()
	globals.Config = loadConfig()
	globals.Font = newFont()
	globals.G = loadGameData()
}
