package oss

import (
	"fmt"
	"testing"
	"time"
)

func TestUpload(t *testing.T) {
	os := NewQiniu("HrtyrRh6okam2Kh3tjhAiduw1gKsmh4tsRLmGgk6",
		"nlYVi-buu8dMbZnr7ECcKsMWZOYcdD_trgc6NCXw",
		"south", time.Hour)
	//err := os.UploadFile("/Users/bighuangbee/Pictures/01古龙峡/IMG_5940_sd.JPG", "")

	for i := 0; i < 100; i++ {
		fmt.Println(os.Sign())
		time.Sleep(time.Second)
	}
}
