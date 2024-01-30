package images

import (
	_ "embed"
)

var (
	//go:embed tanks.png
	Tanks_png []byte

	//go:embed terrains.png
	Terrains_png []byte
)
