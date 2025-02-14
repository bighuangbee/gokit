package main

/*
#cgo LDFLAGS: -framework Security
*/
import "C"
import (
	"fmt"
	os2 "github.com/bighuangbee/gokit/storage/oss"
)

func main() {
	os := os2.NewQiniu("HrtyrRh6okam2Kh3tjhAiduw1gKsmh4tsRLmGgk6", "nlYVi-buu8dMbZnr7ECcKsMWZOYcdD_trgc6NCXw", "south")
	err := os.UploadFile("/Users/bighuangbee/code/face/video_log/output/detect/0101/2号机/144616.000_164.mp4/best_1_73_c-0.2610770_q-0.9971491.jpg", "")
	//fmt.Println(os.Sign())

	//err := os.Download("effect/0212/effect_best_1_73_c-0.2610770_q-0.9971491.mp4", "./1.mp4")
	//err := os.Download("1739515632028_319313834567733959_黄伟华.jpg", "./1.jpg")
	fmt.Println(err)
}
