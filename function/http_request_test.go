package function

import (
	"fmt"
	"testing"
	"time"
)

func TestFrom(t *testing.T) {
	addr := "http://192.168.100.203:6002"
	filename := "/Users/bighuangbee/code/face/face_recognize/recognize/test_data/DSC08085.JPG"

	fmt.Println("注册人脸", "filename", filename)

	t1 := time.Now()
	respBody, err := HttpFormPost(addr+"/face/registe", map[string]string{
		"video_filename": filename,
	}, &FormFile{
		FormField: "file",
		Filename:  filename,
	})

	fmt.Println("face register request", err, string(respBody))
	fmt.Println("耗时", time.Now().Sub(t1).Milliseconds(), "ms")
}
