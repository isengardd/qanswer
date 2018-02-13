package qanswer

import (
	"qanswer/config"
	"qanswer/ocr"
	"qanswer/proto"
)

//Ocr ocr 识别图片文字
type Ocr interface {
	GetText(imgPath string) (string, error)
}

//NewOcr 使用哪种ocr识别
func NewOcr(cfg *config.Config) Ocr {
	if cfg.OcrType == proto.OcrTesseract {
		return ocr.NewTesseract(cfg)
	}
	return ocr.NewBaidu(cfg)
}
