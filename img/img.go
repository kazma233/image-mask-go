package img

import (
	"image"
	"image-mask/utils"
	"image/color"
	"os"

	"github.com/disintegration/imaging"
	"go.uber.org/zap"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type (
	imgFont string
)

const (
	// Regular 常规
	Regular imgFont = "REGULAR"
	// Medium 中等
	Medium imgFont = "MEDIUM"
	// Bold 粗体
	Bold imgFont = "BOLD"

	regularPath = "resource/font/SourceHanSansCN-Regular.otf"
	mediumPath  = "resource/font/SourceHanSansCN-Medium.otf"
	boldPath    = "resource/font/SourceHanSansCN-Bold.otf"
)

var (
	log       = utils.GetLog()
	fontStyle = map[imgFont]*opentype.Font{}
)

func init() {
	fontStyle[Regular] = getFontPanic(regularPath)
	fontStyle[Medium] = getFontPanic(mediumPath)
	fontStyle[Bold] = getFontPanic(boldPath)
	log.Debug("init font info: ", zap.Any("val", fontStyle))
}

// WriteColorMask 在图片上填满色块
func WriteColorMask(bgImg image.Image, width, hight int, c color.Color, p image.Point) image.Image {
	return imaging.Paste(
		bgImg,
		imaging.New(width, hight, c),
		p,
	)
}

// WriteWordMask 写文字水印
func WriteWordMask(bgImg image.Image, word string, imgFont imgFont, c color.Color, fontSize, dpi float64, pt fixed.Point26_6) (image.Image, fixed.Int26_6, error) {
	face, err := opentype.NewFace(fontStyle[imgFont], &opentype.FaceOptions{
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

func getFontPanic(path string) *opentype.Font {
	fontFile, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		log.Error("open font file err: ", zap.Error(err))
		return nil
	}
	font, err := opentype.ParseReaderAt(fontFile)
	if err != nil {
		log.Error("parse font file err: ", zap.Error(err))
		panic(err)
	}

	return font
}
