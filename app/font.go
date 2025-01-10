package main

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	//go:embed assets/RobotoMono-Regular.ttf
	robotoMonoTtf []byte

	//go:embed assets/RobotoMono-Bold.ttf
	robotoMonoBoldTtf []byte

	fontFaceSource *text.GoTextFaceSource
	FontFace22     *text.GoTextFace
	FontFace18     *text.GoTextFace
	FontFace16     *text.GoTextFace
	FontFace14     *text.GoTextFace

	fontFaceBoldSource *text.GoTextFaceSource
	FontFaceBold16     *text.GoTextFace
	FontFaceBold18     *text.GoTextFace
)

func LoadFonts() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(robotoMonoTtf))
	if err != nil {
		log.Fatal(err)
	}

	sBold, err := text.NewGoTextFaceSource(bytes.NewReader(robotoMonoBoldTtf))
	if err != nil {
		log.Fatal(err)
	}

	fontFaceSource = s
	fontFaceBoldSource = sBold

	FontFace22 = &text.GoTextFace{
		Source: fontFaceSource,
		Size:   22,
	}
	FontFace18 = &text.GoTextFace{
		Source: fontFaceSource,
		Size:   18,
	}
	FontFace16 = &text.GoTextFace{
		Source: fontFaceSource,
		Size:   16,
	}
	FontFace14 = &text.GoTextFace{
		Source: fontFaceSource,
		Size:   14,
	}
	FontFaceBold16 = &text.GoTextFace{
		Source: fontFaceBoldSource,
		Size:   16,
	}
	FontFaceBold18 = &text.GoTextFace{
		Source: fontFaceBoldSource,
		Size:   18,
	}
}
