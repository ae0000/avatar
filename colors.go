package avatar

import (
	"errors"
	"image"
	"image/color"
)

// Colors for background
var (
	Red      = image.Uniform{color.RGBA{230, 25, 75, 255}}
	Green    = image.Uniform{color.RGBA{60, 180, 75, 255}}
	Yellow   = image.Uniform{color.RGBA{255, 225, 25, 255}}
	Blue     = image.Uniform{color.RGBA{0, 130, 200, 255}}
	Orange   = image.Uniform{color.RGBA{245, 130, 48, 255}}
	Purple   = image.Uniform{color.RGBA{145, 30, 180, 255}}
	Cyan     = image.Uniform{color.RGBA{70, 240, 240, 255}}
	Magenta  = image.Uniform{color.RGBA{240, 50, 230, 255}}
	Lime     = image.Uniform{color.RGBA{210, 245, 60, 255}}
	Pink     = image.Uniform{color.RGBA{250, 190, 190, 255}}
	Teal     = image.Uniform{color.RGBA{0, 128, 128, 255}}
	Lavender = image.Uniform{color.RGBA{230, 190, 255, 255}}
	Brown    = image.Uniform{color.RGBA{170, 110, 40, 255}}
	Beige    = image.Uniform{color.RGBA{255, 250, 200, 255}}
	Maroon   = image.Uniform{color.RGBA{128, 0, 0, 255}}
	Mint     = image.Uniform{color.RGBA{170, 255, 195, 255}}
	Olive    = image.Uniform{color.RGBA{128, 128, 0, 255}}
	Coral    = image.Uniform{color.RGBA{255, 215, 180, 255}}
	Navy     = image.Uniform{color.RGBA{0, 0, 128, 255}}
	Grey     = image.Uniform{color.RGBA{128, 128, 128, 255}}
	Gold     = image.Uniform{color.RGBA{251, 184, 41, 255}}
)

var errInvalidFormat = errors.New("invalid format")

// parseHexColorFast was found here:
// https://stackoverflow.com/questions/54197913/parse-hex-string-to-image-color
func parseHexColorFast(s string) (c color.RGBA, err error) {
	c.A = 0xff

	if s[0] != '#' {
		return c, errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = errInvalidFormat
	}
	return
}

// TODO add some more colors
func defaultColor(initial string) image.Uniform {
	switch initial {
	case "A", "0":
		return Red
	case "B", "1":
		return Green
	case "C", "2":
		return Yellow
	case "D", "3":
		return Blue
	case "E", "4":
		return Orange
	case "F", "5":
		return Purple
	case "G", "6":
		return Lime
	case "H", "7":
		return Magenta
	case "I", "8":
		return Pink
	case "J", "9":
		return Cyan
	case "K":
		return Teal
	case "L":
		return Lavender
	case "M":
		return Brown
	case "N":
		return Beige
	case "O":
		return Maroon
	case "P":
		return Mint
	case "Q":
		return Olive
	case "R":
		return Coral
	case "S":
		return Navy
	case "T":
		return Gold
	default:
		return Grey
	}
}
