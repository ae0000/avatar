package avatar

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
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

// SetFontFacePath sets the font to do the business with
func SetFontFacePath(f string) {
	fontFacePath = f
}

// var sourceDir string

// func init() {
// 	// We need to set the source directory for the font
// 	_, filename, _, ok := runtime.Caller(0)
// 	if !ok {
// 		panic("No caller information")
// 	}
// 	sourceDir = path.Dir(filename)
// }

// ToDisk saves the image to disk
func ToDisk(initials, path string) {
	rgba, err := createAvatar(initials)
	if err != nil {
		log.Println(err)
		return
	}

	// Save image to disk
	out, err := os.Create(path)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer out.Close()

	b := bufio.NewWriter(out)

	err = png.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// ToHTTP sends the image to a http.ResponseWriter (as a PNG)
func ToHTTP(initials string, w http.ResponseWriter) {
	rgba, err := createAvatar(initials)
	if err != nil {
		log.Println(err)
		return
	}

	b := new(bytes.Buffer)
	key := fmt.Sprintf("avatar%s", initials) // for Etag

	err = png.Encode(b, rgba)
	if err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(b.Bytes())))
	w.Header().Set("Cache-Control", "max-age=2592000") // 30 days
	w.Header().Set("Etag", `"`+key+`"`)

	if _, err := w.Write(b.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
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
	fontBytes, err := ioutil.ReadFile(fontPath) //fmt.Sprintf("%s/%s", sourceDir, fontFaceName))
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

func createAvatar(initials string) (*image.RGBA, error) {
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
	background := defaultColor(text[0:1])
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
