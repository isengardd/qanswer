package screenshot

import (
	"image"
	"log"
	"os/exec"
	"strings"
	"time"

	"qanswer/config"
	"qanswer/proto"
	"qanswer/util"
)

//Android android
type Android struct{}

//NewAndroid new
func NewAndroid(cfg *config.Config) *Android {
	return new(Android)
}

//GetImage 通过adb获取截图
func (android *Android) GetImage() (img image.Image, err error) {

	//	t1 := time.Now().Year()   //年
	//	t2 := time.Now().Month()  //月
	//	t3 := time.Now().Day()    //日
	//	t4 := time.Now().Hour()   //小时
	//	t5 := time.Now().Minute() //分钟
	//	t6 := time.Now().Second() //秒
	//

	datestring := strings.Replace(time.Now().Format("2006-01-02 15:04:05"), ":", "", -1)
	datestring = strings.Replace(datestring, "-", "", -1)
	datestring = strings.Replace(datestring, " ", "", -1)
	targetImagePath := "/sdcard/screenshot" + datestring + ".png"
	if config.GetConfig().Debug {
		log.Printf("图片保存路径:%v", targetImagePath)
	}

	err = exec.Command("adb", "shell", "screencap", "-p", targetImagePath).Run()
	if err != nil {
		return
	}
	originImagePath := proto.ImagePath + "origin.png"
	err = exec.Command("adb", "pull", targetImagePath, originImagePath).Run()
	if err != nil {
		return
	}

	//	err = exec.Command("adb", "shell", "rm", targetImagePath).Run()
	//	if err != nil {
	//		return
	//	}

	img, err = util.OpenPNG(originImagePath)
	return
}
