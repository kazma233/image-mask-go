package resource

import (
	"embed"
	"image-mask/utils"

	"golang.org/x/image/font/opentype"
)

const (
	// Regular 常规
	Regular string = "font/SourceHanSansCN-Regular.otf"
)

var (
	log = utils.GetSugaredLogger()
	//go:embed font/*
	fontFS embed.FS
)

// GetFontByte 获取字体数据
func GetFontByte(fn string) []byte {
	fbs, err := fontFS.ReadFile(fn)
	if err != nil {
		log.Errorf("read font file error: %v", err)
		return nil
	}

	return fbs
}

// GetFont 获取字体对象
func GetFont(fn string) *opentype.Font {
	fbs := GetFontByte(fn)
	if fbs == nil || len(fbs) <= 0 {
		return nil
	}

	font, err := opentype.Parse(fbs)
	if err != nil {
		log.Errorf("parse font file err: %v", err)
		return nil
	}

	return font
}
