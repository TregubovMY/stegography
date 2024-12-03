package lsb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"io"

	"github.com/TregubovMY/stegography/bitmanip"
	"github.com/TregubovMY/stegography/utils"
)

// Decode выполняет стеганографическое декодирование считывателя
// с ранее закодированными данными с помощью функции Encode и записывает результат в io.Writer.

//export Decode
func Decode(c []byte) ([]byte, error) {
	carrier := io.NopCloser(bytes.NewReader(c))
	RGBAImage, _, err := utils.GetImageAsRGBA(carrier)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга контейнера данных: %v", err)
	}

	width := RGBAImage.Bounds().Dx()
	height := RGBAImage.Bounds().Dy()

	dataBytes := make([]byte, 0, 2048)
	resultBytes := make([]byte, 0, 2048)

	dataCount := extractDataCount(RGBAImage)

	var count int

	for x := 0; x < width && dataCount > 0; x++ {
		for y := 0; y < height && dataCount > 0; y++ {
			if count >= dataSizeHeaderReservedBytes {
				c := RGBAImage.RGBAAt(x, y)
				dataBytes = append(dataBytes,
					bitmanip.GetLastTwoBits(c.R),
					bitmanip.GetLastTwoBits(c.G),
					bitmanip.GetLastTwoBits(c.B),
				)
				dataCount -= 3
			} else {
				count += 4
			}
		}
	}

	if dataCount < 0 {
		//remove bytes that are not part of data and mistakenly added
		dataBytes = dataBytes[:len(dataBytes)+dataCount]
	}

	dataBytes = align(dataBytes) // len(dataBytes) must be aliquot of 4

	for i := 0; i < len(dataBytes); i += 4 {
		resultBytes = append(resultBytes, bitmanip.ConstructByteOfQuartersAsSlice(dataBytes[i:i+4]))
	}

	// if _, err = result.Write(resultBytes); err != nil {
	// 	return nil, err
	// }

	return resultBytes, nil
}

func align(dataBytes []byte) []byte {
	switch len(dataBytes) % 4 {
	case 1:
		dataBytes = append(dataBytes, byte(0), byte(0), byte(0))
	case 2:
		dataBytes = append(dataBytes, byte(0), byte(0))
	case 3:
		dataBytes = append(dataBytes, byte(0))
	}
	return dataBytes
}

func extractDataCount(RGBAImage *image.RGBA) int {
	dataCountBytes := make([]byte, 0, 16)

	width := RGBAImage.Bounds().Dx()
	height := RGBAImage.Bounds().Dy()

	count := 0

	for x := 0; x < width && count < dataSizeHeaderReservedBytes; x++ {
		for y := 0; y < height && count < dataSizeHeaderReservedBytes; y++ {
			pixel := RGBAImage.RGBAAt(x, y)
			dataCountBytes = append(dataCountBytes,
				bitmanip.GetLastTwoBits(pixel.R),
				bitmanip.GetLastTwoBits(pixel.G),
				bitmanip.GetLastTwoBits(pixel.B),
			)
			count += 4
		}
	}

	dataCountBytes = append(dataCountBytes, byte(0))

	var bs = []byte{
		bitmanip.ConstructByteOfQuartersAsSlice(dataCountBytes[:4]),
		bitmanip.ConstructByteOfQuartersAsSlice(dataCountBytes[4:8]),
		bitmanip.ConstructByteOfQuartersAsSlice(dataCountBytes[8:12]),
		bitmanip.ConstructByteOfQuartersAsSlice(dataCountBytes[12:]),
	}

	return int(binary.LittleEndian.Uint32(bs))
}

// func Sum(a, b int) int { return a + b }

// func MyFoo(data []byte) int { return len(data) }

// func MyFooErr(data []byte) (int, error) { return len(data), nil }
