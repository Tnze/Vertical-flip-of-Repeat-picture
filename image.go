package main

import (
	"fmt"
	"github.com/Tnze/CoolQ-Golang-SDK/v2/cqp"
	"github.com/Tnze/CoolQ-Golang-SDK/v2/cqp/util"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"
	"path/filepath"
)

// procImage 处理一张图片消息并返回处理结果
// panic if msg is not a image like [CQ:image,file=1.jpg]
func procImage(msg string) (string, error) {
	imgDir := imageReg.FindStringSubmatch(msg)[1]
	imgDir = cqp.GetImage(imgDir)

	// 检查处理后图片是否已存在
	baseProcImg, procImgDir := procImgPath(imgDir, "png")
	if _, err := os.Stat(procImgDir); err != nil && os.IsNotExist(err) {
		err := convertImg(procImgDir, imgDir)
		if err != nil {
			return "", err
		}
	}

	return util.CQCode("image", "file", baseProcImg), nil
}

func procImgPath(name, targetExt string) (string, string) {
	baseName := filepath.Base(name)
	baseWithoutExt := baseName[:len(baseName)-len(filepath.Ext(baseName))]
	base := baseWithoutExt + "." + targetExt
	return filepath.Join(cqp.AppID, base), filepath.Join(imgFold, base)
}

func convertImg(dst, src string) error {
	// 读取源图片
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("无法打开图片，%v", err)
	}
	defer f.Close()
	img, typ, err := image.Decode(f)
	if err != nil {
		return fmt.Errorf("无法解析图片，%v", err)
	}
	Debugf("图片处理", "成功解析一张%s图片", typ)

	// 垂直翻转
	img = vFlip(img)

	// 写入新图片
	df, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("无法创建图片文件，%v", err)
	}
	defer df.Close()
	if err := png.Encode(df, img); err != nil {
		return fmt.Errorf("无法编码图片，%v", err)
	}
	return nil
}

// 上下翻转
func vFlip(m image.Image) image.Image {
	mb := m.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, mb.Dx(), mb.Dy()))
	for x := mb.Min.X; x < mb.Max.X; x++ {
		for y := mb.Min.Y; y < mb.Max.Y; y++ {
			//  设置像素点
			dst.Set(x, mb.Max.Y-y, m.At(x, y))
		}
	}
	return dst
}
