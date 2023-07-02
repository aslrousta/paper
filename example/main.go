package main

import (
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"

	"github.com/aslrousta/paper"
)

func main() {
	in, err := os.Open("girl.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	girl, err := jpeg.Decode(in)
	if err != nil {
		log.Fatal(err)
	}

	if err := render(girl, "nostalgia", paper.Nostalgia); err != nil {
		log.Fatal(err)
	}
	if err := render(girl, "sepia", paper.Sepia); err != nil {
		log.Fatal(err)
	}
	if err := render(girl, "night", paper.Night); err != nil {
		log.Fatal(err)
	}
}

func render(im image.Image, name string, theme paper.Theme) error {
	p := paper.New(theme, im.Bounds().Dx(), im.Bounds().Dy())
	draw.Draw(p, im.Bounds(), im, image.Point{}, draw.Src)

	out, err := os.Create(name + ".jpg")
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, p, nil)
}
