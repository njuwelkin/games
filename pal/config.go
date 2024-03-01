package main

import "github.com/njuwelkin/games/pal/mkf"

type Config struct {
	GamePath   string
	SavePath   string
	ShaderPath string
	WordLength mkf.DWORD
}

func loadConfig() Config {
	return Config{
		GamePath:   "./",
		SavePath:   "./",
		ShaderPath: "./",
		WordLength: 10,
	}
}
