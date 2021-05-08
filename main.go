package main

import (
	"bytes"
	"image"
	"image-mask/entity"
	"image-mask/resource"
	"image-mask/utils"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

var (
	// log
	log = utils.GetSugaredLogger()
)

func main() {
	router := gin.Default()

	router.Static("/assets", "./resource")
	router.StaticFile("/favicon.ico", "./resource/favicon.ico")

	router.GET("/", func(c *gin.Context) {
		c.File("./resource/html/index.html")
	})

	router.GET("/img", func(c *gin.Context) {
		imageInfo := &entity.ImageInfo{}
		if err := c.BindQuery(&imageInfo); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		log.Info("image info %s", imageInfo)

		img, err := processImage(imageInfo)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		buf := new(bytes.Buffer)
		_ = jpeg.Encode(buf, img, nil)
		_, _ = c.Writer.Write(buf.Bytes())
	})

	log.Errorf("服务运行：%s", router.Run(":9000"))
}

func processImage(imageInfo *entity.ImageInfo) (img image.Image, err error) {
	// open bg
	img, err = imaging.Open(resource.Root + resource.BG)
	if err != nil {
		return
	}

	// open chara
	charaImg, err := imaging.Open(resource.Root + resource.Chara1)
	if err != nil {
		return
	}
	img = imaging.Overlay(img, charaImg, image.Pt(imageInfo.CharaX, imageInfo.CharaY), 1)

	// Write Font
	var position int = 100
	for _, r := range "中文 English" {
		fontFace, err := utils.PreWordMask(
			entity.WordMaskPreInfo{
				Word: string(r),
				Font: resource.GetFont(resource.Regular),
				Size: 65,
				Dpi:  100,
			},
		)

		if err != nil {
			log.Errorf("写出文字失败: %v", err)
			break
		}

		_bgImg, err := utils.WriteWordMask(
			fontFace,
			entity.WordMaskInfo{
				BgImg: img,
				Word:  string(r),
				ColorPoint: entity.ColorPoint{
					C: color.RGBA{R: utils.UColor(), G: utils.UColor(), B: utils.UColor(), A: 255},
					X: 250,
					Y: position,
				},
			},
		)
		if err != nil {
			log.Errorf("写出文字失败: %s", err)
			break
		}
		position += 100
		img = _bgImg
	}

	// Write color
	img = utils.WriteColorMask(
		entity.WriteColorInfo{
			BgImage: img,
			ColorBoxInfo: entity.ColorBox{
				Width: 100,
				High:  100,
				Point: image.Pt(imageInfo.BoxX, imageInfo.BoxY),
				Color: color.RGBA{R: utils.UColor(), G: utils.UColor(), B: utils.UColor(), A: 255},
			},
		},
	)

	return
}
