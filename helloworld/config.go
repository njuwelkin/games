package main

import (
	"encoding/json"
	"image/color"
	"log"
	"os"
)

type Config struct {
	ScreenWidth     int        `json:"screenWidth"`
	ScreenHeight    int        `json:"screenHeight"`
	Title           string     `json:"title"`
	BgColor         color.RGBA `json:"bgColor"`
	ShipSpeedFactor float64    `json:"shipSpeedFactor"`
}

func newDefaultConfig() *Config {
	return &Config{
		ScreenWidth:  640,
		ScreenHeight: 480,
		Title:        "外星人入侵",
		BgColor:      color.RGBA{R: 230, G: 230, B: 230, A: 255},
	}
}

func loadConfig() *Config {
	f, err := os.Open("./config.json")
	if err != nil {
		return newDefaultConfig()
	}

	var cfg Config
	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		log.Fatalf("json.Decode failed: %v\n", err)
	}

	return &cfg
}
