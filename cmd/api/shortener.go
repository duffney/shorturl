package main

import (
	"math"
	"strings"
	"time"
)

type Shorten struct {
	Long_url  string    `json:"long_url"`
	Short_url string    `json:"short_url"`
	CreatedAt time.Time `json:"-"`
}

const base62Digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const shortenerAddress = "http://localhost:4000/v1/" // #TODO: Move to config

func DecimalToBase62(n int64) string {
	if n == 0 {
		return "0"
	}

	base62 := make([]byte, 0)
	radix := int64(62)

	for n > 0 {
		remainder := n % radix
		base62 = append([]byte{base62Digits[remainder]}, base62...)
		n /= radix
	}

	return string(base62)
}

func Base62ToDecimal(s string) int64 {
	var decimalNumber int64

	for i, c := range s {
		decimalNumber += int64(strings.IndexByte(base62Digits, byte(c))) * int64(math.Pow(62, float64(len(s)-i-1)))
	}

	return decimalNumber
}

// func main() {
// 	decimalNumber := int64(2009215674938)
// 	base62Representation := DecimalToBase62(decimalNumber)
// 	fmt.Printf("Base62 representation of %d is %s\n", decimalNumber, base62Representation)
// }

// Context: https://chat.openai.com/share/19b461d8-90cc-42df-8985-9335caa81a9d
