package utils

import (
	"fmt"
	"github.com/TregubovMY/stegography/bitmanip"
	"image"
	"image/draw"
	"io"
)

// ReadData читает данные из reader и отправляет их в канал bytes по четвертям
// байта (2 бита). Если происходит ошибка, то она отправляется в канал errChan.
func ReadData(reader io.Reader, bytes chan<- byte, errChan chan<- error) {
	b := make([]byte, 1)
	for {
		if _, err := reader.Read(b); err != nil {
			if err == io.EOF {
				break
			}
			errChan <- fmt.Errorf("ошибка чтения данных %v", err)
			return
		}
		for _, b := range bitmanip.QuartersOfByte(b[0]) {
			bytes <- b
		}
	}
	close(bytes)
}

// GetImageAsRGBA декодирует изображение из reader, копирует его в новый
// *image.RGBA, возвращая его, формат изображения и ошибку, если она
// возникла.
func GetImageAsRGBA(reader io.Reader) (*image.RGBA, string, error) {
	img, format, err := image.Decode(reader)
	if err != nil {
		return nil, format, fmt.Errorf("ошибка декодирования изображения в rgba: %v", err)
	}

	RGBAImage := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(RGBAImage, RGBAImage.Bounds(), img, img.Bounds().Min, draw.Src)

	return RGBAImage, format, nil
}

// SetColorSegment записывает байт, полученный из data, в последние два
// бита colorSegment. Если data закрыт, то функция возвращает false, nil.
// Если errChan содержит ошибку, то функция возвращает false, ошибку.

func SetColorSegment(colorSegment *byte, data <-chan byte, errChan <-chan error) (hasMoreBytes bool, err error) {
	select {
	case byte, ok := <-data:
		if !ok {
			return false, nil
		}
		*colorSegment = bitmanip.SetLastTwoBits(*colorSegment, byte)
		return true, nil

	case err := <-errChan:
		return false, err
	}
}
