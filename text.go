package main

import (
	"bytes"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	_ "embed"
)

var (
	//go:embed assets/fonts/Moralerspace_v1.0.0/MoralerspaceNeon-Regular.ttf
	moralerspaceNeonRegular []byte
	faceSource              *text.GoTextFaceSource
	faces                   = make(map[float64]*text.GoTextFace)
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(moralerspaceNeonRegular))
	if err != nil {
		log.Fatal(err)
	}
	faceSource = s

	var normalFace = &text.GoTextFace{
		Source: faceSource,
		Size:   12,
	}
	faces[12] = normalFace
}

func Write(screen *ebiten.Image, t string, x, y, size float64, clr ebiten.ColorScale) {
	face := faces[size]
	if face == nil {
		face = &text.GoTextFace{
			Source: faceSource,
			Size:   size,
		}
		faces[size] = face
	}
	var op = text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale = clr
	op.LayoutOptions = text.LayoutOptions{
		LineSpacing:    1,
		PrimaryAlign:   text.AlignCenter,
		SecondaryAlign: text.AlignCenter,
	}
	text.Draw(screen, t, face, &op)
}

func WriteLeft(screen *ebiten.Image, t string, x, y, size float64, clr ebiten.ColorScale) {
	face := faces[size]
	if face == nil {
		face = &text.GoTextFace{
			Source: faceSource,
			Size:   size,
		}
		faces[size] = face
	}
	var op = text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale = clr
	op.LayoutOptions = text.LayoutOptions{
		LineSpacing:    1,
		PrimaryAlign:   text.AlignStart,
		SecondaryAlign: text.AlignCenter,
	}
	text.Draw(screen, t, face, &op)
}
