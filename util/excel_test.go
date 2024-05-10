package util

import (
	"log"
	"testing"
)

func Test_getExcelFile(t *testing.T) {
	name := getExcelFile()
	log.Println(name)
}
