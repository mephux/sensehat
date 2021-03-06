package sensehat

import (
	"errors"
	"os"
)

func (d *Device) SetPixels(list PixelList) error {

	if len(list) != 64 {
		return errors.New("pixel list has the wrong len")
	}

	f, err := os.OpenFile(d.FDPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		return err
	}

	defer f.Close()

	for i, pixel := range list {
		f.Seek(int64(d.Map[d.Rotation][i/8][i%8]*2), 0)
		f.Write(pixel.Pack())
	}

	return nil
}

func (d *Device) GetPixels() (PixelList, error) {
	var list PixelList

	for x := 0; x <= 7; x++ {
		for y := 0; y <= 7; y++ {
			p, err := d.GetPixel(x, y)

			if err != nil {
				return list, err
			}

			list = append(list, p)
		}
	}

	return list, nil
}

func (d *Device) SetPixel(x, y int, pixel *Pixel) (*Pixel, error) {
	if x > 7 || x < 0 {
		return pixel, errors.New("X position must be between 0 and 7")
	}

	if y > 7 || y < 0 {
		return pixel, errors.New("Y position must be between 0 and 7")
	}

	if !pixel.Valid() {
		return pixel, errors.New("invalid pixel")
	}

	f, err := os.OpenFile(d.FDPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		return pixel, err
	}

	defer f.Close()

	f.Seek(int64(d.Map[d.Rotation][x][y]*2), 0)

	_, err = f.Write(pixel.Pack())

	return pixel, err
}

func (d *Device) GetPixel(x, y int) (*Pixel, error) {
	pixel := &Pixel{}

	if x > 7 || x < 0 {
		return pixel, errors.New("X position must be between 0 and 7")
	}

	if y > 7 || y < 0 {
		return pixel, errors.New("Y position must be between 0 and 7")
	}

	f, err := os.Open(d.FDPath)

	if err != nil {
		return pixel, err
	}

	defer f.Close()

	f.Seek(int64(d.Map[d.Rotation][x][y]*1), 0)

	buffer := make([]byte, 2)

	if _, err := f.Read(buffer); err != nil {
		return pixel, err
	}

	pixel.Unpack(buffer)

	return pixel, nil
}

func (d *Device) SetAllPixels(pixel *Pixel) {
	for x := 0; x <= 7; x++ {
		for y := 0; y <= 7; y++ {
			d.SetPixel(x, y, pixel)
		}
	}
}

func (d *Device) ClearScreen() {
	d.SetAllPixels(&Pixel{
		R: 0,
		G: 0,
		B: 0,
	})
}
