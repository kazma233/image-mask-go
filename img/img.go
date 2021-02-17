package img

import (
	"image"
	"image-mask/utils"
	"image/color"

	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var (
	log = utils.GetSugaredLogger()
)

// WriteColorMask 在图片上填满色块
func WriteColorMask(bgImg image.Image, width, hight int, c color.Color, p image.Point) image.Image {
	return imaging.Paste(
		bgImg,
		imaging.New(width, hight, c),
		p,
	)
}

// WriteWordMask 写文字水印
func WriteWordMask(bgImg image.Image, word string, ft *opentype.Font, c color.Color, fontSize, dpi float64, pt fixed.Point26_6) (image.Image, fixed.Int26_6, error) {
	face, err := opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, 0, err
	}

	dstImg := image.NewRGBA(bgImg.Bounds())
	drawer := font.Drawer{
		Dst:  dstImg,
		Src:  image.NewUniform(c),
		Face: face,
		Dot:  pt,
	}

	drawer.DrawString(word)

	return imaging.OverlayCenter(bgImg, dstImg, 1), drawer.MeasureString(word), nil
}
