package game

import "github.com/njuwelkin/games/pal/pkg/mkf"

type Config struct {
	GamePath   string
	SavePath   string
	ShaderPath string
	WordLength mkf.DWORD
}

func loadConfig() Config {
	return Config{
		GamePath:   "./data",
		SavePath:   "./",
		ShaderPath: "./",
		WordLength: 10,
	}
}
