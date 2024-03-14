package DrawImg

import (
	"Click_CAPTCHA_Go/Env"
	"Click_CAPTCHA_Go/Utils"
	Config "Click_CAPTCHA_Go/config"
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"image/png"
	"log"
	"math"
	"math/rand"
	"strings"
)

func DrawCaptchaImg(DrawMode int) (imgByte []byte, CoordinateSlice []Coordinates, Answer []AnswerImg) {
	imageWidth := 600
	imageHeight := 400
	context := gg.NewContext(imageWidth, imageHeight)
	context.DrawImage(utils.BackgroundByteSlice[rand.Intn(len(utils.BackgroundByteSlice))], 0, 0)
	var Coordinate Coordinates
	var CoordinateSliceTemp []Coordinates
	var AnswerImgTemp []AnswerImg
	switch DrawMode {
	case Env.ImagesClickMode:

		IconSlice := SelectImage(utils.IconByteSlice, 3)
		for i := 0; i < 3; i++ {
			Coordinate.X = rand.Intn(472) + 64
			Coordinate.Y = rand.Intn(272) + 64
			fmt.Println(Coordinate.X, Coordinate.Y)
			CoordinateSliceTemp = append(CoordinateSliceTemp, Coordinate)

			var IconBuf bytes.Buffer
			err := png.Encode(&IconBuf, IconSlice[i])
			if err != nil {
				log.Println(err.Error())
				return nil, nil, nil
			}
			AnswerImgTemp = append(AnswerImgTemp, AnswerImg{ImgByte: IconBuf.Bytes()})
			context.DrawImage(IconSlice[i], Coordinate.X, Coordinate.Y)
		}
	case Env.WordsClickMode:
		Text := strings.TrimSpace(SelectText(utils.TextSlice))
		for _, value := range Text {
			var Coordinate Coordinates
			for {
				Coordinate.X = rand.Intn(472) + 64
				Coordinate.Y = rand.Intn(272) + 64
				if !IsOverlap(Coordinate, CoordinateSliceTemp) {
					break
				}
			}
			fmt.Println(Coordinate.X, Coordinate.Y)
			fmt.Println(string(value))
			CoordinateSliceTemp = append(CoordinateSliceTemp, Coordinate)
			context.SetFontFace(utils.FontSlice[rand.Intn(len(utils.FontSlice))])
			context.SetColor(utils.GenerateRandomColor())
			//context.Rotate(gg.Radians(rand.Float64()*30 - 15))
			context.DrawString(string(value), float64(Coordinate.X), float64(Coordinate.Y))
		}

	}

	var buf bytes.Buffer
	err := png.Encode(&buf, context.Image())
	if err != nil {
		log.Println(err.Error())
		return nil, nil, nil
	}
	return buf.Bytes(), CoordinateSliceTemp, AnswerImgTemp
}

// SelectImage 从传入的图片切片中随机选择x张图片(Fisher-Yates洗牌算法)
func SelectImage(ImageArray []image.Image, x int) []image.Image {
	selected := make([]image.Image, x)
	for i := 0; i < x; i++ {
		index := rand.Intn(len(ImageArray)-i) + i
		selected[i] = ImageArray[index]
		ImageArray[i], ImageArray[index] = ImageArray[index], ImageArray[i]
	}
	return selected
}

// SelectText 从传入的文字切片中随机选择一个文字
func SelectText(TextArray []string) string {
	return TextArray[rand.Intn(len(TextArray))]
}

func IsOverlap(newCoordinate Coordinates, existingCoordinates []Coordinates) bool {
	for _, coordinate := range existingCoordinates {
		// 这里我们假设如果两个字符的位置在30个像素内，那么它们就会重叠
		if math.Abs(float64(newCoordinate.X-coordinate.X)) < Config.GlobalConfig.DrawWordSize && math.Abs(float64(newCoordinate.Y-coordinate.Y)) < Config.GlobalConfig.DrawWordSize {
			return true
		}
	}
	return false
}
