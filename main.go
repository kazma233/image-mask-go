package main

import (
	"image"
	"image-mask/img"
	"image-mask/resource"
	"image-mask/utils"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/disintegration/imaging"
	"github.com/flopp/go-findfont"
	"go.uber.org/zap"
	"golang.org/x/image/math/fixed"
)

const (
	// size
	baseWidth  = 400
	baseHeight = 200
	fileWidth  = 800
	fileHeight = 400
)

var (
	// log
	log = utils.GetSugaredLogger()
	// ui
	a               fyne.App
	baseWindows     fyne.Window
	imgModalWindows fyne.Window
	showLabel       *widget.Label
	// position
	colorX int
	colorY int
	// chara position
	charaX int
	charaY int
	//
	imgCanvas *canvas.Image
)

func init() {
	for _, fn := range findfont.List() {
		lowerFn := strings.ToLower(fn)
		if strings.Contains(lowerFn, "pingfang") || strings.Contains(lowerFn, "msyh") {
			log.Debugf("find font: %s", fn)
			err := os.Setenv("FYNE_FONT", fn)
			if err != nil {
				log.Errorf("set env FYNE_FONT error: %v", err)
			}

			break
		}
	}
}

func main() {
	a = app.New()

	baseWindows = a.NewWindow("图片处理")
	baseWindows.Resize(fyne.Size{Width: baseWidth, Height: baseHeight})
	baseWindows.CenterOnScreen()

	showLabel = widget.NewLabel("just a demo")

	// 输入框
	colorXEntry := widget.NewEntryWithData(binding.IntToString(binding.BindInt(&colorX)))
	colorXEntry.SetPlaceHolder("色块x")
	colorYEntry := widget.NewEntryWithData(binding.IntToString(binding.BindInt(&colorY)))
	colorYEntry.SetPlaceHolder("色块y")

	// chara输入框
	charaXEntry := widget.NewEntryWithData(binding.IntToString(binding.BindInt(&charaX)))
	charaXEntry.SetPlaceHolder("角色x")
	charaYEntry := widget.NewEntryWithData(binding.IntToString(binding.BindInt(&charaY)))
	charaYEntry.SetPlaceHolder("角色y")

	baseWindows.SetContent(container.NewVBox(
		showLabel,
		container.NewGridWithColumns(
			3, widget.NewLabel("输入色块的x,y: "),
			colorXEntry, colorYEntry,
		),
		container.NewGridWithColumns(
			3, widget.NewLabel("输入角色的x,y: "),
			charaXEntry, charaYEntry,
		),
		widget.NewButton("选取excel", selectExcel),
		widget.NewButton("打开预览", processImage),
	))

	baseWindows.SetCloseIntercept(func() {
		shutdown()
	})

	baseWindows.Show()
	a.Run()
}

func shutdown() {
	os.Unsetenv("FYNE_FONT")
	closeWindows(imgModalWindows)
	closeWindows(baseWindows)
	a.Quit()
}

func closeWindows(w fyne.Window) {
	if w != nil {
		w.Close()
	}
}

func selectExcel() {
	size := fyne.NewSize(fileWidth, fileHeight)
	fileModalWindows := a.NewWindow("文件选择器")
	fileModalWindows.Resize(size)
	fileModalWindows.CenterOnScreen()
	fileModalWindows.Show()

	fileDialg := dialog.NewFileOpen(func(u fyne.URIReadCloser, e error) {
		fileModalWindows.Hide()
		if e != nil {
			log.Error(e.Error())
		} else if u != nil {
			showLabel.SetText(u.URI().Path())
		}
	}, fileModalWindows)
	fileDialg.Resize(size)
	fileDialg.Show()
}

func processImage() {
	// open bg
	bgImg, err := imaging.Open(resource.Root + resource.BG)
	if err != nil {
		panic(err)
	}

	// open chara
	charaImg, err := imaging.Open(resource.Root + resource.Chara1)
	if err != nil {
		panic(err)
	}
	bgImg = imaging.Overlay(bgImg, charaImg, image.Pt(charaX, charaY), 1)

	// Write Font
	position := 100
	for _, r := range "我怎么知道ABC" {
		_bgImg, size, err := img.WriteWordMask(
			bgImg,
			string(r), resource.GetFont(resource.Regular),
			color.RGBA{R: utils.UColor(), G: utils.UColor(), B: utils.UColor(), A: 255},
			65, 100, fixed.P(250, position),
		)
		if err != nil {
			log.Error("写出文字失败: ", zap.Error(err))
		}
		position += size.Floor()
		bgImg = _bgImg
	}

	// Write color
	bgImg = img.WriteColorMask(bgImg,
		100, 100,
		color.RGBA{R: utils.UColor(), G: utils.UColor(), B: utils.UColor(), A: 255},
		image.Pt(colorX, colorY),
	)

	if imgCanvas == nil {
		imgCanvas = canvas.NewImageFromImage(bgImg)
		imgCanvas.FillMode = canvas.ImageFillOriginal
	}

	if imgModalWindows == nil {
		imgModalWindows = a.NewWindow("图片查看器")
		imgModalWindows.Resize(imgCanvas.Size())
		imgModalWindows.SetContent(imgCanvas)
		imgModalWindows.CenterOnScreen()
		imgModalWindows.Show()
		imgModalWindows.SetCloseIntercept(func() {
			imgModalWindows.Hide()
		})
	} else {
		imgModalWindows.Show()
	}

	imgCanvas.Image = bgImg
	imgCanvas.Refresh()
}
