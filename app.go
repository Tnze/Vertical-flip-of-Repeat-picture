package main

import (
	"fmt"
	"github.com/Tnze/CoolQ-Golang-SDK/v2/cqp"
	"os"
	"path/filepath"
	"regexp"
)

//go:generate cqcfg -c .
// cqp: 名称: 复读图片之垂直翻转
// cqp: 版本: 1.1.1:1
// cqp: 作者: Tnze
// cqp: 简介: 当一张图在群内被重复两次以上时，将该图片上下翻转并发送，以有效打断复读
func main() { cqp.Main() }

func init() {
	cqp.AppID = "online.jdao.VerticalFlipOfRepeatPicture"
	cqp.Enable = onEnable
	cqp.GroupMsg = onGroupMsg

}

var imgFold string

// 当插件启用
func onEnable() int32 {
	imgFold = filepath.Join("data", "image", cqp.AppID)
	err := os.Mkdir(imgFold, 0644)
	if err != nil && !os.IsExist(err) {
		Errorf("初始化", "无法创建用于发送图片的文件夹，%v", err)
		return -1
	}
	return 0
}

// 当收到群消息
func onGroupMsg(subType, msgID int32, fromGroup, fromQQ int64, fromAnonymous, msg string, font int32) int32 {
	g := getGroup(fromGroup)
	g.Lock()
	defer g.Unlock()
	if g.lastMsg != msg {
		g.repeatCount = 1
		g.lastMsg = msg
	} else {
		g.repeatCount++
		if g.repeatCount == 2 && isImage(g.lastMsg) {
			Infof("复读", "发现一次复读，%s", msg)
			img, err := procImage(msg)
			if err != nil {
				Errorf("复读", "图片处理失败，%v", err)
				return Ignore
			}
			cqp.SendGroupMsg(fromGroup, img)
		}
	}
	return Ignore
}

var imageReg = regexp.MustCompile(`\[CQ:image,file=([^"]*)\]`)

func isImage(msg string) bool {
	// 必须要这条消息包含且仅包含一张图片
	// 才认为这是一条图片消息
	return msg != "" && msg == imageReg.FindString(msg)
}

// Errorf格式化输出错误日志
func Errorf(tp, format string, args ...interface{}) {
	cqp.AddLog(cqp.Error, tp, fmt.Sprintf(format, args...))
}

// Infof格式化输出信息日志
func Infof(tp, format string, args ...interface{}) {
	cqp.AddLog(cqp.Info, tp, fmt.Sprintf(format, args...))
}

// Debugf格式化输出调试日志
func Debugf(tp, format string, args ...interface{}) {
	cqp.AddLog(cqp.Debug, tp, fmt.Sprintf(format, args...))
}

const (
	Ignore    int32 = 0 //忽略消息
	Intercept       = 1 //拦截消息
)
