// Пакет предоставляет утилиты для манипулирования битами одного байта
package bitmanip

import (
	"encoding/binary"
)

const (
	firstQuarter  = 192 // 1100 0000
	secondQuarter = 48  // 0011 0000
	thirdQuarter  = 12  // 0000 1100
	fourthQuarter = 3   // 0000 0011
)

// QuartersOfByte возвращает четыре четверти бита в байте
func QuartersOfByte(b byte) [4]byte {
	return [4]byte{
		b & firstQuarter >> 6,
		b & secondQuarter >> 4,
		b & thirdQuarter >> 2,
		b & fourthQuarter,
	}
}

func clearLastTwoBits(b byte) byte {
	return b & byte(252) // 1111 1100
}

// SetLastTwoBits изменяет последние два бита байта на значение value.
func SetLastTwoBits(b byte, value byte) byte {
	return clearLastTwoBits(b) | value
}

// GetLastTwoBits возвращает последние два бита байта.
func GetLastTwoBits(b byte) byte {
	return b & fourthQuarter
}

// ConstructByteOfQuarters создаёт байт из четырёх двухбитных частей.
func ConstructByteOfQuarters(first, second, third, fourth byte) byte {
	return (((first << 6) | (second << 4)) | third<<2) | fourth
}

// ConstructByteOfQuartersAsSlice создаёт байт из массива четырёх двухбитных значений.
func ConstructByteOfQuartersAsSlice(b []byte) byte {
	return ConstructByteOfQuarters(b[0], b[1], b[2], b[3])
}

// QuartersOfBytesOf преобразует заданное число uint32 в 16 двухбитовых значений (4 байта).
func QuartersOfBytesOf(counter uint32) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, counter)
	quarters := make([]byte, 16)
	for i := 0; i < 16; i += 4 {
		quartersOfByte := QuartersOfByte(bs[i/4])
		quarters[i] = quartersOfByte[0]
		quarters[i+1] = quartersOfByte[1]
		quarters[i+2] = quartersOfByte[2]
		quarters[i+3] = quartersOfByte[3]
	}

	return quarters
}
