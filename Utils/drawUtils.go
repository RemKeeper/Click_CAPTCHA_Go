package utils

import (
	"Click_CAPTCHA_Go/Env"
	Config "Click_CAPTCHA_Go/config"
	"bytes"
	"context"
	"fmt"
	"github.com/beevik/etree"
	"github.com/golang/freetype/truetype"
	"github.com/kanrichan/resvg-go"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
)

var (
	BackgroundByteSlice []image.Image
	IconByteSlice       []image.Image
	TextSlice           []string
	FontSlice           []font.Face
)

func ReadBackgroundsToByteSlice() {
	dir, err := os.Open(Env.BackgroundPath)
	if err != nil {
		log.Println(err.Error(), " 背景目录未找到或无权限")
		return
	}
	FileNames, err := dir.Readdirnames(-1)
	if err != nil {
		log.Println(err.Error(), " 读取背景文件资源失败")
		return
	}

	// 创建一个channel来收集结果
	results := make(chan image.Image, len(FileNames))

	// 创建一个WaitGroup来等待所有goroutine完成
	var wg sync.WaitGroup

	for _, FileName := range FileNames {
		// 为每个文件创建一个goroutine
		wg.Add(1)
		go func(FileName string) {
			defer wg.Done()
			file, err := os.Open(Env.BackgroundPath + FileName)
			if err != nil {
				return
			}
			decode, _, err := image.Decode(file)
			if err != nil {
				return
			}
			results <- resize.Resize(600, 400, decode, resize.Lanczos3)
		}(FileName)
	}

	// 创建一个新的goroutine来关闭结果channel，当所有读取goroutine完成后
	go func() {
		wg.Wait()
		close(results)
	}()

	// 从结果channel中收集结果
	for result := range results {
		BackgroundByteSlice = append(BackgroundByteSlice, result)
	}
}

func ReadIconsToByteSlice() {
	dir, err := os.Open(Env.IconPath)
	if err != nil {
		log.Println(err.Error(), "图标目录未找到或无权限")
		return
	}
	FileNames, err := dir.Readdirnames(-1)
	if err != nil {
		return
	}

	// 创建一个channel来收集结果
	results := make(chan image.Image, len(FileNames))

	// 创建一个WaitGroup来等待所有goroutine完成
	var wg sync.WaitGroup

	for _, File := range FileNames {
		// 为每个文件创建一个goroutine
		wg.Add(1)
		go func(File string) {
			defer wg.Done()
			SvgIco := ReadSvg(Env.IconPath + File)
			if !isImageZeroValue(SvgIco) {
				results <- SvgIco
			}
		}(File)
	}

	// 创建一个新的goroutine来关闭结果channel，当所有读取goroutine完成后
	go func() {
		wg.Wait()
		close(results)
	}()

	// 从结果channel中收集结果
	for result := range results {
		IconByteSlice = append(IconByteSlice, result)
	}
}

func ReadWordTextToSlice() {
	dir, err := os.Open(Env.TextPath)
	if err != nil {
		log.Println(err.Error(), "文本目录未找到或无权限")
		return
	}
	FIleNames, err := dir.Readdirnames(-1)
	if err != nil {
		return
	}
	for _, File := range FIleNames {
		file, err := os.ReadFile(Env.TextPath + File)
		if err != nil {
			return
		}
		TextSlice = append(TextSlice, strings.Split(string(file), "\n")...)
	}
}

func ReadWordFontToByteSlice() {
	dir, err := os.Open(Env.FontPath)
	if err != nil {
		log.Println(err.Error(), "字体目录未找到或无权限")
		return
	}
	FIleNames, err := dir.Readdirnames(-1)
	if err != nil {
		return
	}
	for _, File := range FIleNames {
		file, err := os.ReadFile(Env.FontPath + File)
		if err != nil {
			return
		}
		parseFont, err := truetype.Parse(file)
		if err != nil {
			return
		}
		face := truetype.NewFace(parseFont, &truetype.Options{
			Size: Config.GlobalConfig.DrawWordSize,
		})
		FontSlice = append(FontSlice, face)
	}
}

func generateRandomHexColor() string {
	red := rand.Intn(156) + 100
	green := rand.Intn(156) + 100
	blue := rand.Intn(156) + 100
	sprintf := fmt.Sprintf("#%02X%02X%02X", red, green, blue)
	return sprintf
}

func GenerateRandomColor() color.Color {
	red := rand.Intn(56) + 200
	green := rand.Intn(56) + 200
	blue := rand.Intn(56) + 200
	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: 255}
}

func ReadSvg(fileName string) image.Image {
	fmt.Println(fileName)
	doc := etree.NewDocument()
	err := doc.ReadFromFile(fileName)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	root := doc.SelectElement("svg")
	for _, element := range root.SelectElements("path") {
		for i, attr := range element.Attr {
			if attr.Key == "stroke" || attr.Key == "fill" {
				element.Attr[i].Value = generateRandomHexColor()
			}
		}
	}
	SvgBytes, err := doc.WriteToBytes()
	if err != nil {
		log.Println(err)
		return nil
	}

	worker, err := resvg.NewDefaultWorker(context.Background())
	defer worker.Close()
	if err != nil {
		log.Println(err)
		return nil
	}
	render, err := worker.Render(SvgBytes)
	if err != nil {
		log.Println(err)
		return nil
	}
	decode, _, err := image.Decode(bytes.NewReader(render))
	if err != nil {
		log.Println(err)
		return nil
	}
	return resize.Resize(64, 64, decode, resize.Lanczos3)
}

func isImageZeroValue(img image.Image) bool {
	bounds := img.Bounds()
	// 遍历图像的每个像素
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			// 如果有任何一个像素的 RGBA 值不全为零，则不是全零值
			if r != 0 || g != 0 || b != 0 || a != 0 {
				return false
			}
		}
	}
	return true
}
