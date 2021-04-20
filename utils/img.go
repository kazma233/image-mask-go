package utils

import (
	"image"
	"image-mask/entity"

	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var (
	log = GetSugaredLogger()
)

// WriteColorMask 在图片上填满色块
func WriteColorMask(info entity.WriteColorInfo) image.Image {
	return imaging.Paste(
		info.BgImage,
		imaging.New(info.ColorBoxInfo.Width, info.ColorBoxInfo.Hight, info.ColorBoxInfo.Color),
		info.ColorBoxInfo.Point,
	)
}

// PreWordMask 获取文字水印
func PreWordMask(info entity.WordMaskPreInfo) (font.Drawer, error) {
	face, err := opentype.NewFace(info.Font, &opentype.FaceOptions{
		Size:    info.Size,
		DPI:     info.Dpi,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return font.Drawer{}, err
	}

	drawer := font.Drawer{
		Face: face,
	}

	return drawer, nil
}

// WriteWordMask 写文字水印
func WriteWordMask(info entity.WordMaskInfo) (image.Image, error) {
	bgImg := info.BgImg
	drawer := info.Drawer

	dstImg := image.NewRGBA(bgImg.Bounds())
	drawer.Dst = dstImg
	drawer.Dot = info.Pt
	drawer.Src = image.NewUniform(info.C)

	drawer.DrawString(info.Word)

	return imaging.OverlayCenter(bgImg, dstImg, 1), nil
}

// WriteFontCenter 写出文字（中间）
func WriteFontCenter(info entity.WordMaskCenterInfo) (image.Image, error) {
	drawer, err := PreWordMask(
		entity.WordMaskPreInfo{
			Word: info.Word,
			Font: info.Font,
			Size: info.Size,
			Dpi:  info.Dpi,
		},
	)
	fSize := drawer.MeasureString(info.Word)
	if err != nil {
		return nil, err
	}

	log.Debugf("fSize: %v", fSize)

	w := (info.Width - fSize.Floor()) / 2
	wordMaskImage, err := WriteWordMask(
		entity.WordMaskInfo{
			BgImg:  info.BgImg,
			Drawer: drawer,
			Word:   info.Word,
			ColorPoint: entity.ColorPoint{
				C:  info.C,
				Pt: fixed.P(w, info.Y),
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return wordMaskImage, nil
}
