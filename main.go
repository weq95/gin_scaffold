package main

import (
	"flag"
	"fmt"
	"github.com/gin_scaffiold/common/lib"
	"github.com/gin_scaffiold/router"
	"github.com/skip2/go-qrcode"
	"image/color"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	endpoint = flag.String("endpoint", "", "input endpoint dashboard or server")
	config   = flag.String("config", "", "input config file like ./conf/dev/")
)

func main() {
	flag.Parse()
	if *endpoint == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *config == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *endpoint == "dashboard" {
		_ = lib.InitModule("./conf/dev/")

		defer lib.Destroy()

		router.HttpServerRun()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		router.HttpServerStop()
	} else {
		_ = lib.InitModule("./conf/dev/")

		defer lib.Destroy()

		router.HttpServerRun()

		fmt.Println("start Server:代理服务器已启动")
		//todo 添加代理服务器

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit
	}

}

func init() {
	qr, err := qrcode.New("http://www.flysnow.org/", qrcode.Medium)
	if err != nil {
		return
	}

	//生成目录
	dir := "./static/qrcode/" + time.Now().Format("200601") + "/"
	fileName := time.Now().Format("20060102150405") + ".png"
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	qr.BackgroundColor = color.White
	//255,215,0
	qr.ForegroundColor = color.RGBA{R: 255, G: 215, A: 255}
	err = qr.WriteFile(256, dir+fileName)
	if err != nil {
		panic(err)
	}
}
