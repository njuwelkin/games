package main

type global struct {
	Text   TextLib
	Config Config
	Font   Font
}

var globalSetting global

func initGlobalSetting() {
	globalSetting.Text = loadText()
	globalSetting.Config = loadConfig()
	globalSetting.Font = newFont()
}
