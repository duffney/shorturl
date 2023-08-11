package main

import (
	"fmt"
	"net/url"
)

func main() {
	short_url := "http://localhost:4000/v1/7VQTFRr8s"

	fmt.Println(url.Parse(short_url))
}
