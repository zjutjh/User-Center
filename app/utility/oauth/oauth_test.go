package oauth

import (
	"log"
	"testing"
)

func TestHttp(t *testing.T) {
	c, err := CheckByOauth("", "")
	log.Println(c)
	log.Println(err)
}
