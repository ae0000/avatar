// Package avatar generates colorful avatar images from initials.
// It creates 500x500px PNG images with customizable colors and built-in caching.
package avatar

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

const (
	defaultfontFace = "Roboto-Bold.ttf" //SourceSansVariable-Roman.ttf"
	fontSize        = 210.0
	imageWidth      = 500.0
	imageHeight     = 500.0
	dpi             = 72.0
	spacer          = 20
	textY           = 320
)

var fontFacePath = ""

// SetFontFacePath sets the path to the TrueType font file to use for rendering avatar text.
// If not set, the default font "Roboto-Bold.ttf" will be used.
func SetFontFacePath(f string) {
	fontFacePath = f
}

// ToDisk generates an avatar image from the given initials and saves it as a PNG file
// at the specified path. The initials can be a full name (e.g., "John Doe" becomes "JD")
// or explicit initials (e.g., "AE"). Colors are automatically assigned based on the first character.
// Returns an error if avatar generation or file writing fails.
func ToDisk(initials, path string) error {
	return saveToDisk(initials, path, "", "")
}

// ToDiskCustom generates an avatar image with custom colors and saves it as a PNG file.
// The bgColor and fontColor parameters should be hex color strings (e.g., "#FF0000").
// Empty strings use default colors. Returns an error if avatar generation or file writing fails.
func ToDiskCustom(initials, path, bgColor, fontColor string) error {
	return saveToDisk(initials, path, bgColor, fontColor)
}

// saveToDisk saves the image to disk
func saveToDisk(initials, path, bgColor, fontColor string) error {
	rgba, err := createAvatar(initials, bgColor, fontColor)
	if err != nil {
		return fmt.Errorf("create avatar: %w", err)
	}

	// Save image to disk
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer out.Close()

	b := bufio.NewWriter(out)

	err = png.Encode(b, rgba)
	if err != nil {
		return fmt.Errorf("encode png: %w", err)
	}

	err = b.Flush()
	if err != nil {
		return fmt.Errorf("flush buffer: %w", err)
	}

	return nil
}

// ToHTTP generates an avatar image from the given initials and writes it as a PNG
// to the provided http.ResponseWriter. Automatically sets appropriate headers including
// Content-Type, Cache-Control (30 days), and ETag for browser caching.
// Returns an error if avatar generation or writing fails.
func ToHTTP(initials string, w http.ResponseWriter) error {
	return saveToHTTP(initials, "", "", w)
}

// ToHTTPCustom generates an avatar image with custom colors and writes it as a PNG
// to the provided http.ResponseWriter. The bgColor and fontColor parameters should be
// hex color strings (e.g., "#FF0000"). Empty strings use default colors.
// Returns an error if avatar generation or writing fails.
func ToHTTPCustom(initials, bgColor, fontColor string, w http.ResponseWriter) error {
	return saveToHTTP(initials, bgColor, fontColor, w)
}

// saveToHTTP sends the image to a http.ResponseWriter (as a PNG)
func saveToHTTP(initials, bgColor, fontColor string, w http.ResponseWriter) error {
	rgba, err := createAvatar(initials, bgColor, fontColor)
	if err != nil {
		return fmt.Errorf("create avatar: %w", err)
	}

	b := new(bytes.Buffer)
	key := fmt.Sprintf("avatar%s%s%s", initials, bgColor, fontColor) // for Etag

	err = png.Encode(b, rgba)
	if err != nil {
		return fmt.Errorf("encode png: %w", err)
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(b.Bytes())))
	w.Header().Set("Cache-Control", "max-age=2592000") // 30 days
	w.Header().Set("Etag", `"`+key+`"`)

	if _, err := w.Write(b.Bytes()); err != nil {
		return fmt.Errorf("write image: %w", err)
	}

	return nil
}

func cleanString(incoming string) string {
	incoming = strings.TrimSpace(incoming)

	// If its something like "firstname surname" get the initials out
	split := strings.Split(incoming, " ")
	if len(split) == 2 {
		incoming = split[0][0:1] + split[1][0:1]
	}

	// Max length of 2
	if len(incoming) > 2 {
		incoming = incoming[0:2]
	}

	// To upper and trimmed
	return strings.ToUpper(strings.TrimSpace(incoming))
}

func getFont(fontPath string) (*truetype.Font, error) {
	if fontPath == "" {
		fontPath = defaultfontFace
	}
	// Read the font data.
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}

	return freetype.ParseFont(fontBytes)
}

var imageCache sync.Map

func getImage(initials string) *image.RGBA {
	value, ok := imageCache.Load(initials)

	if !ok {
		return nil
	}

	image, ok2 := value.(*image.RGBA)
	if !ok2 {
		return nil
	}
	return image
}

func setImage(initials string, image *image.RGBA) {
	imageCache.Store(initials, image)
}

func createAvatar(initials, bgColor, fontColor string) (*image.RGBA, error) {
	// Make sure the string is OK
	text := cleanString(initials)

	// Check cache
	cachedImage := getImage(text)
	if cachedImage != nil {
		return cachedImage, nil
	}

	// Load and get the font
	f, err := getFont(fontFacePath)
	if err != nil {
		return nil, err
	}

	// Setup the colors, text white, background based on first initial
	textColor := image.White
	if fontColor != "" {
		c, err := parseHexColorFast(fontColor)
		if err == nil {
			textColor = &image.Uniform{c}
		}
	}
	background := defaultColor(text[0:1])
	if bgColor != "" {
		c, err := parseHexColorFast(bgColor)
		if err == nil {
			background = image.Uniform{c}
		}
	}

	rgba := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	draw.Draw(rgba, rgba.Bounds(), &background, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(fontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(textColor)
	c.SetHinting(font.HintingFull)

	// We need to convert the font into a "font.Face" so we can read the glyph
	// info
	to := truetype.Options{}
	to.Size = fontSize
	face := truetype.NewFace(f, &to)

	// Calculate the widths and print to image
	xPoints := []int{0, 0}
	textWidths := []int{0, 0}

	// Get the widths of the text characters
	for i, char := range text {
		width, ok := face.GlyphAdvance(rune(char))
		if !ok {
			return nil, err
		}

		textWidths[i] = int(float64(width) / 64)
	}

	// TODO need some tests for this
	if len(textWidths) == 1 {
		textWidths[1] = 0
	}

	// Get the combined width of the characters
	combinedWidth := textWidths[0] + spacer + textWidths[1]

	// Draw first character
	xPoints[0] = int((imageWidth - combinedWidth) / 2)
	xPoints[1] = int(xPoints[0] + textWidths[0] + spacer)

	for i, char := range text {
		pt := freetype.Pt(xPoints[i], textY)
		_, err := c.DrawString(string(char), pt)
		if err != nil {
			return nil, err
		}
	}

	// Cache it
	setImage(text, rgba)

	return rgba, nil
}
