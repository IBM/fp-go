package codec

import (
	"log"
	"testing"
)

func TestStringCoded(t *testing.T) {

	sType := String()

	res := sType.Decode(10)

	log.Println(res)
}
