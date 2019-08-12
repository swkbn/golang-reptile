package main

import (
	"github.com/kataras/iris"
	"log"
	"net/http"
	"xuandan/controller"
	"xuandan/thecrawler"
)

func main() {

	go thecrawler.AllPage()
	go thecrawler.PartPage(1, 10)
	go thecrawler.OverdueGoods()
	app := iris.New()
	app.Get("/api/goodsdetale", controller.HandlerGoodsDetale)
	app.Get("/api/ofungirdled", controller.HandlerOfungirdled)
	app.Get("/api/figurewith", controller.HandlerFigurewith)
	app.Get("/api/videothe", controller.HandlerVideoThe)
	app.Get("/api/goods", controller.Handler)
	app.Get("/api/goods/searchItemId", controller.HandlerItemId)
	app.Get("/api/goods/search", controller.HandlerSearch)
	app.Get("/api/goods/category", controller.HandlerSearchCategory)
	//创建监听
	err := app.Run(iris.Server(&http.Server{Addr: ":9090"}))
	//l, _ := net.Listen("tcp", "127.0.0.1:9090")
	//err := app.Run(iris.Listener(l))
	if err != nil {
		log.Println("创建失败：", err)
		return
	}

}
