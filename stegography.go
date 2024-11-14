package main

import (
	"bytes"
	"fmt"
	"github.com/TregubovMY/stegography/stegify_methods/lsb"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"
)

func main() {
	// Указываем пути к файлам
	carrierFile := "input.png"
	dataFile := "data.txt"
	encodedFile := "encoded.png"
	decodedDataFile := "decoded.txt"

	if err := createTestImage(carrierFile); err != nil {
		fmt.Println("Ошибка при создании тестового изображения:", err)
		return
	}

	carrier, err := os.Open(carrierFile)
	if err != nil {
		log.Fatalf("Ошибка при открытии контейнера %s: %v", carrierFile, err)
	}
	defer carrier.Close()

	data, err := os.Open(dataFile)
	if err != nil {
		log.Fatalf("Ошибка при открытии файла данных %s: %v", dataFile, err)
	}
	defer data.Close()

	encoded, err := os.Create(encodedFile)
	if err != nil {
		log.Fatalf("Ошибка при создании закодированного файла %s: %v", encodedFile, err)
	}
	defer encoded.Close()

	err = lsb.Encode(carrier, data, encoded)
	if err != nil {
		log.Fatalf("Ошибка при кодировании данных: %v", err)
	}
	fmt.Println("Данные успешно закодированы в изображение:", encodedFile)

	encodedImage, err := os.Open(encodedFile)
	if err != nil {
		log.Fatalf("Ошибка при открытии закодированного файла %s: %v", encodedFile, err)
	}
	defer encodedImage.Close()

	// Буфер для хранения раскодированных данных
	decodedData := new(bytes.Buffer)

	err = lsb.Decode(encodedImage, decodedData)
	if err != nil {
		log.Fatalf("Ошибка при декодировании данных: %v", err)
	}
	fmt.Println("Данные успешно декодированы из изображения")

	err = os.WriteFile(decodedDataFile, decodedData.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Ошибка при сохранении декодированных данных: %v", err)
	}
	fmt.Printf("Декодированные данные сохранены в файл: %s\n", decodedDataFile)

	fmt.Println("Декодированные данные:")
	io.Copy(os.Stdout, decodedData)
}

// Генерация PNG изображения с красным фоном
func createTestImage(fileName string) error {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255}) // красный фон
		}
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
