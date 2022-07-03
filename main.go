package main

import (
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	uuid "github.com/satori/go.uuid"
	"github.com/skip2/go-qrcode"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"
	_ "time"
)

func main() {
	tableNumber := "Стол №10"
	width := 900
	height := 1020
	uri := "https://www.example.ru/my-venue/42"

	resultImageName := buildResultImageName(uri)
	qrCodeFilename := createQRCode(width, uri)
	backgroundFilename := createBackground(tableNumber, width, height)
	createResultImage(qrCodeFilename, backgroundFilename, height-width, resultImageName)
}

func buildResultImageName(uri string) string {
	uriSplit := strings.Split(uri, "/")
	return fmt.Sprintf("qr-code-%s-%s.png", uriSplit[len(uriSplit)-2], uriSplit[len(uriSplit)-1])
}

func createResultImage(qrCodeFilename string, backgroundFilename string, point int, resultImageName string) {
	backgroundImageStream, err := os.Open(backgroundFilename)
	checkError(err)
	backgroundImage, err := png.Decode(backgroundImageStream)
	checkError(err)
	defer backgroundImageStream.Close()
	qrCodeStream, err := os.Open(qrCodeFilename)
	checkError(err)
	qrImage, err := png.Decode(qrCodeStream)
	checkError(err)
	defer qrCodeStream.Close()

	offset := image.Pt(0, point)
	bounds := backgroundImage.Bounds()

	resultImage := image.NewRGBA(bounds)
	draw.Draw(resultImage, bounds, backgroundImage, image.ZP, draw.Src)
	draw.Draw(resultImage, qrImage.Bounds().Add(offset), qrImage, image.ZP, draw.Over)
	resultStream, err := os.Create(resultImageName)
	checkError(err)
	png.Encode(resultStream, resultImage)
	defer resultStream.Close()
}

func createQRCode(width int, uri string) string {
	uuidQRCode, err := uuid.NewV4()
	checkError(err)
	qrCodeFilename := fmt.Sprintf("qr-code-%s.png", uuidQRCode)

	err = qrcode.WriteFile(
		uri,
		qrcode.Medium,
		width,
		qrCodeFilename,
	)
	fmt.Println(err)
	return qrCodeFilename
}

func createBackground(tableNumber string, width int, height int) string {
	uuidBackground, err := uuid.NewV4()
	checkError(err)
	backgroundFilename := fmt.Sprintf("background-%s.png", uuidBackground)
	background := buildBackground(width, height, tableNumber)
	file, err := os.Create(backgroundFilename)
	checkError(err)
	err = png.Encode(file, background)
	checkError(err)

	return backgroundFilename
}

func buildBackground(width int, height int, initials string) *image.RGBA {
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	background := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(background, background.Bounds(), &image.Uniform{C: white}, image.Point{}, draw.Src)
	drawText(background, initials)
	return background
}

func drawText(canvas *image.RGBA, text string) error {
	var (
		fgColor  image.Image
		fontFace *truetype.Font
		err      error
		fontSize = 100.0
	)
	fgColor = image.Black
	fontFace, err = freetype.ParseFont(goregular.TTF)
	fontDrawer := &font.Drawer{
		Dst: canvas,
		Src: fgColor,
		Face: truetype.NewFace(fontFace, &truetype.Options{
			Size:    fontSize,
			Hinting: font.HintingFull,
		}),
	}

	fontDrawer.Dot = fixed.Point26_6{X: 4000, Y: 7500}
	fontDrawer.DrawString(text)
	return err
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
