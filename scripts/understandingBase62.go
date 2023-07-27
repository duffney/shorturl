package main

import (
	"fmt"
	"math"
	"strings"
)

func main() {
	const base62Digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	n := int64(123)
	radix := int64(62)
	base62 := make([]byte, 0)
	remainder := n % radix
	fmt.Println(remainder)
	base62Digit := []byte{base62Digits[remainder]}
	/* In Go, a string is a sequence of bytes, so base62Digits[remainder] returns a
	byte value. The byte value is then used to create a byte slice with a
	single element.
	*/
	fmt.Println(string(base62Digit))
	base62 = append(base62Digit, base62...)
	n /= radix
	fmt.Println(n)
	remainder = n % radix
	fmt.Println(remainder)
	base62 = append([]byte{base62Digits[remainder]}, base62...)
	n /= radix
	fmt.Println(n)
	fmt.Println(string(base62))
	// convert back

	//var decimalNumber int64
	for i, c := range base62 {
		base62ToDecmial := struct {
			Character     string
			IndexPosition int64
			PowerOf       int64
			DecimalValue  int64
		}{
			Character:     string(c),
			IndexPosition: int64(strings.IndexByte(base62Digits, byte(c))),
			PowerOf:       int64(math.Pow(62, float64(len(base62)-i-1))),
			DecimalValue:  int64(strings.IndexByte(base62Digits, byte(c))) * int64(math.Pow(62, float64(len(base62)-i-1))),
		}
		fmt.Println("len", len(base62))
		fmt.Println("count", i)
		fmt.Println("secondNumber", len(base62)-i-1)
		fmt.Println(base62ToDecmial)

		//base62IndexPosition := int64(strings.IndexByte(base62Digits, byte(c)))
		//fmt.Println("base62Character", string(c))
		//fmt.Println("base62IndexPosition", base62IndexPosition)
		//powerOf := int64(math.Pow(62, float64(len(base62)-i-1)))
		//fmt.Println("Powerof", powerOf)
	}
}
