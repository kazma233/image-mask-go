package resource

import "embed"

const (
	// BG 背景图
	BG = "img/bg.png"
	// Chara1 人物图片
	Chara1 = "img/img_chara_1.png"
)

// ImgFS image files
//go:embed img/*
var ImgFS embed.FS
