package main

import gopaint "github.com/vitorduarte/go-paint"

func main() {
	painter, err := gopaint.NewPainterFromConfigFile("config.json")
	if err != nil {
		panic(err)
	}

	img, err := painter.Paint()
	if err != nil {
		panic(err)
	}

	gopaint.SaveImage(img, "output.png")
}
