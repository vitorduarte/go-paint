package gopaint

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
)

type Painter struct {
	Config PaintConfig
}

type PaintConfig struct {
	Template string       `json:"template"`
	Masks    []MaskConfig `json:"masks"`
}

type MaskConfig struct {
	Path  string `json:"path"`
	Color string `json:"color"`
}

func NewPainterFromConfigFile(configPath string) (p Painter, err error) {
	config, err := ReadConfigFile(configPath)
	if err != nil {
		return
	}
	return Painter{Config: config}, nil
}

func ReadConfigFile(path string) (config PaintConfig, err error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return
	}

	return
}

func (p *Painter) Paint() (canva *image.RGBA, err error) {
	canva, err = openTemplateImage(p.Config.Template)
	if err != nil {
		err = fmt.Errorf("Could not open template image: %v", err)
		return
	}

	for _, mask := range p.Config.Masks {
		var fill_color color.RGBA
		if mask.Color == "random" {
			fill_color = generateRandomColor()
		} else {
			fill_color, err = hexToRGBA(mask.Color)
			if err != nil {
				err = fmt.Errorf("Could not parse color: %v", err)
				return
			}
		}

		err = fillImageWithColorWithMask(canva, mask.Path, fill_color)
		if err != nil {
			err = fmt.Errorf("Could not parse color: %v", err)
			return
		}
	}
	return
}

func openTemplateImage(path string) (canva *image.RGBA, err error) {
	templateImage, err := openImage(path)
	if err != nil {
		return
	}

	canva = image.NewRGBA(templateImage.Bounds())
	draw.Draw(canva, canva.Bounds(), templateImage, image.ZP, draw.Src)
	return
}

func openImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func generateRandomColor() color.RGBA {
	return color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255}
}

func hexToRGBA(hex string) (c color.RGBA, err error) {
	// Consider alpha to be 255 if it was not specified
	c.A = 0xff

	// Format hexadecimal color according the number of digits provided
	switch len(hex) {
	case 9:
		_, err = fmt.Sscanf(hex, "#%02x%02x%02x%02x", &c.R, &c.G, &c.B, &c.A)
	case 7:
		_, err = fmt.Sscanf(hex, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(hex, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("Invalid hex color: %s", hex)
	}

	return
}

func fillImageWithColorWithMask(canva *image.RGBA, maskPath string, color color.RGBA) (err error) {
	fill := &image.Uniform{color}

	mask, err := openImage(maskPath)
	if err != nil {
		return
	}

	draw.DrawMask(canva, canva.Bounds(), fill, image.ZP, mask, image.ZP, draw.Over)
	return
}

func SaveImage(img image.Image, outputPath string) (err error) {
	f, err := os.Create(outputPath)
	if err != nil {
		return
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		return
	}

	return nil
}
