package oss

import (
	"fmt"
	"testing"
)

func TestUpload(t *testing.T) {
	os := NewQiniu("HrtyrRh6okam2Kh3tjhAiduw1gKsmh4tsRLmGgk6", "nlYVi-buu8dMbZnr7ECcKsMWZOYcdD_trgc6NCXw", "south")
	//err := os.UploadFile("/Users/bighuangbee/Pictures/01古龙峡/IMG_5940_sd.JPG", "")
	fmt.Println(os)
}
