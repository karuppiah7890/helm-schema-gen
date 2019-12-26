package main

import (
	"fmt"
	"github.com/karuppiah7890/go-jsonschema-generator"
)

func main() {
	data := map[string]interface{}{
		"something": map[string]interface{}{
			"okay": "dokey",
		},
	}
	s := &jsonschema.Document{}
	s.ReadDeep(&data)
	fmt.Println(s)
}
