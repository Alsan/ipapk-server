package utils

import (
	"image"
	"image/png"
	"os"
)

func SaveIcon(icon image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := png.Encode(file, icon); err != nil {
		return err
	}
	return nil
}
