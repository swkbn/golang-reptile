package controller

import (
	"encoding/json"
	"github.com/kataras/iris"
	"log"
	"strconv"
	"xuandan/models"
	"xuandan/thecrawler"
)

func init() {

}
func HandlerGoodsDetale(ctx iris.Context) {
	ctx.ResponseWriter().Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	ctx.ResponseWriter().Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	ctx.ResponseWriter().Header().Set("content-type", "application/json")             //返回数据格式是json

	Itemid := ctx.FormValue("goodsid")
	db := models.Session
	goodsid, _ := strconv.Atoi(Itemid)
	var goodsItem models.GoodsItem
	var goodsindex int
	db.Model(&models.GoodsItem{}).Where("goods_id = ?", goodsid).Count(&goodsindex)
	if goodsindex == 0 {
		ctx.JSON(map[string]string{"errcode": "数据不存在"})
		return
	}
	db.Model(&models.GoodsItem{}).Where("goods_id = ?", goodsid).Find(&goodsItem)
	goodsdetal := thecrawler.FormatConversion(goodsItem)
	_, err := ctx.JSON(goodsdetal)
	if err != nil {
		log.Println("返回错误：", err)
	}
}

func conversion(fqcat string) int {

	var newFqcat int
	switch fqcat {
	case "":
		newFqcat = 0
	case "全部":
		newFqcat = 0
	case "女装":
		newFqcat = 1
	case "男装":
		newFqcat = 2
	case "内衣":
		newFqcat = 3
	case "美妆":
		newFqcat = 4
	case "配饰":
		newFqcat = 5
	case "鞋品":
		newFqcat = 6
	case "箱包":
		newFqcat = 7
	case "儿童":
		newFqcat = 8
	case "母婴":
		newFqcat = 9
	case "居家":
		newFqcat = 10
	case "美食":
		newFqcat = 11
	case "数码":
		newFqcat = 12
	case "家电":
		newFqcat = 13
	case "车品":
		newFqcat = 15
	case "文体":
		newFqcat = 16
	case "宠物":
		newFqcat = 17
	case "其他":
		newFqcat = 14
	default:
		return 0
	}
	return newFqcat

}

//分类展示
func HandlerSearchCategory(ctx iris.Context) {
	fqcat := ctx.FormValue("fqcat")
	var count []int
	var newFqcat int
	newFqcat = conversion(fqcat)
	db := models.Session
	var goods []thecrawler.Goods
	if newFqcat != 0 {
		db.Model(&models.GoodsItem{}).Order("itemsale2 desc", true).Where("down_type=?", 0).Where("fqcat = ?", newFqcat).Pluck("id", &count)
	} else {
		db.Model(&models.GoodsItem{}).Order("itemsale2 desc", true).Where("down_type=?", 0).Pluck("id", &count)
	}
	for _, val := range count {
		var goodsinform models.GoodsItem
		db.Model(&goodsinform).Where("id=?", val).Find(&goodsinform)
		newGoods := thecrawler.FormatConversion(goodsinform)
		goods = append(goods, newGoods)
	}
	ctx.JSON(goods)

}

//搜索展示
func HandlerSearch(ctx iris.Context) {
	db := models.Session
	var goodsinform []models.GoodsItem

	SearchValue := ctx.FormValue("searchvalue")
	db.Model(&goodsinform).Where("itemtitle LIKE ?", "%"+SearchValue+"%").Order("itemsale2 desc", true).Where("down_type=?", 0).Find(&goodsinform)
	if goodsinform == nil {
		_, err := ctx.JSON(map[string]int{"errcode": 400})
		if err != nil {
			log.Println("写入返回值失败：")
		}
		return
	}
	goods, err := json.Marshal(&goodsinform)
	if err != nil {
		log.Println("json转换失败：", err)
	}
	_, err = ctx.Write(goods)
	if err != nil {
		log.Println("写入返回值失败：", err)
	}
}

//展示
func Handler(ctx iris.Context) {

	ctx.ResponseWriter().Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	ctx.ResponseWriter().Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	ctx.ResponseWriter().Header().Set("content-type", "application/json")             //返回数据格式是json

	var count []int
	log.Println("收到发送的请求")
	var goodsinform []thecrawler.Goods
	printlist := ctx.FormValue("printlist")
	page := ctx.FormValue("page")
	log.Println("请求的参数：", printlist, page)
	var newPrintList string
	switch printlist {
	case "":
		newPrintList = "itemsale2"
	case "两小时销量排行榜":
		newPrintList = "itemsale2"
	case "今日销量排行榜":
		newPrintList = "todaysale"
	case "总销量排行榜":
		newPrintList = "itemsale"
	default:
		ctx.JSON(map[string]int{"erroecode": 400})
		return
	}
	db := models.Session
	db.Model(&models.GoodsItem{}).Order(newPrintList+" desc", true).Where("down_type=?", 0).Pluck("id", &count)
	if count == nil {
		_, err := ctx.JSON(map[string]int{"errcode": 400})
		if err != nil {
			log.Println("写入返回值失败：")
		}
		return
	}
	if page == "" {
		page = "1"
	}
	newPage, _ := strconv.Atoi(page)
	newCount := paging(newPage, count)
	for _, val := range newCount {
		var goods models.GoodsItem
		db.Model(&goods).Where("id= ?", val).Find(&goods)
		newGoods := thecrawler.FormatConversion(goods)
		goodsinform = append(goodsinform, newGoods)
	}
	var item []interface{}
	item = append(item, goodsinform)
	item = append(item, map[string]int{"CounNumb": len(count)})
	_, err := ctx.JSON(item)

	if err != nil {
		log.Println("写入返回值失败：", err)
	}
}

//查询有朋友圈文案的
func HandlerOfungirdled(ctx iris.Context) {
	ctx.ResponseWriter().Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	ctx.ResponseWriter().Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	ctx.ResponseWriter().Header().Set("content-type", "application/json")             //返回数据格式是json

	page := ctx.FormValue("page")
	fqcat := ctx.FormValue("fqcat")
	log.Println("请求来的数据：", page, fqcat)
	var newFqcat int
	newFqcat = conversion(fqcat)
	var count []int
	var goodsinform []thecrawler.Goods
	db := models.Session
	if newFqcat != 0 {
		db.Model(&models.GoodsItem{}).Order("itemsale2 desc", true).Where("down_type=?", 0).Where("Copy_friends_circle_text !=?", "").Where("fqcat = ?", newFqcat).Pluck("id", &count)
	} else {
		db.Model(&models.GoodsItem{}).Order("itemsale2 desc", true).Where("down_type=?", 0).Where("Copy_friends_circle_text !=?", "").Pluck("id", &count)
	}
	if count == nil {
		_, err := ctx.JSON(map[string]int{"errcode": 400})
		if err != nil {
			log.Println("写入返回值失败：")
		}
		return
	}
	if page == "" {
		page = "1"
	}
	newPage, _ := strconv.Atoi(page)

	Newcount := paging(newPage, count)
	for _, val := range Newcount {
		var goods models.GoodsItem
		db.Model(&goods).Where("id= ?", val).Find(&goods)
		newGoods := thecrawler.FormatConversion(goods)
		goodsinform = append(goodsinform, newGoods)
	}
	var item []interface{}
	item = append(item, goodsinform)
	item = append(item, map[string]int{"CounNumb": len(count)})
	_, err := ctx.JSON(item)

	if err != nil {
		log.Println("写入返回值失败：", err)
	}

}

//查有实拍图的
func HandlerFigurewith(ctx iris.Context) {
	ctx.ResponseWriter().Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	ctx.ResponseWriter().Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	ctx.ResponseWriter().Header().Set("content-type", "application/json")             //返回数据格式是json
	fqcat := ctx.FormValue("fqcat")
	page := ctx.FormValue("page")
	log.Println("请求的数据：", fqcat, page)
	var newFqcat int
	newFqcat = conversion(fqcat)
	var count []int
	var goodsinform []thecrawler.Goods
	db := models.Session
	if newFqcat != 0 {
		db.Model(&models.GoodsItem{}).Order("itemsale2 desc", true).Where("down_type=?", 0).Where("Image !=?", "").Where("fqcat = ?", newFqcat).Pluck("id", &count)
	} else {
		db.Model(&models.GoodsItem{}).Order("itemsale2 desc", true).Where("down_type=?", 0).Where("Image !=?", "").Pluck("id", &count)
	}
	if count == nil {
		_, err := ctx.JSON(map[string]int{"errcode": 400})
		if err != nil {
			log.Println("写入返回值失败：")
		}
		return
	}
	if page == "" {
		page = "1"
	}
	newPage, _ := strconv.Atoi(page)
	newCount := paging(newPage, count)
	for _, val := range newCount {
		var goods models.GoodsItem
		db.Model(&goods).Where("id= ?", val).Find(&goods)
		newGoods := thecrawler.FormatConversion(goods)
		goodsinform = append(goodsinform, newGoods)
	}
	var item []interface{}
	item = append(item, goodsinform)
	item = append(item, map[string]int{"CounNumb": len(count)})
	_, err := ctx.JSON(item)

	if err != nil {
		log.Println("写入返回值失败：", err)
	}

}

//查有视频的
func HandlerVideoThe(ctx iris.Context) {
	ctx.ResponseWriter().Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	ctx.ResponseWriter().Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	ctx.ResponseWriter().Header().Set("content-type", "application/json")             //返回数据格式是json

	fqcat := ctx.FormValue("fqcat")
	page := ctx.FormValue("page")
	var newFqcat int
	newFqcat = conversion(fqcat)
	var count []int
	var goodsinform []thecrawler.Goods
	db := models.Session
	if newFqcat != 0 {
		db.Model(&models.GoodsItem{}).Order("itemsale2 desc", true).Where("down_type=?", 0).Where("Main_video_url !=?", "").Where("fqcat = ?", newFqcat).Pluck("id", &count)
	} else {
		db.Model(&models.GoodsItem{}).Order("itemsale2 desc", true).Where("down_type=?", 0).Where("Main_video_url !=?", "").Pluck("id", &count)
	}
	if count == nil {
		_, err := ctx.JSON(map[string]int{"errcode": 400})
		if err != nil {
			log.Println("写入返回值失败：")
		}
		return
	}
	if page == "" {
		page = "1"
	}
	newPage, _ := strconv.Atoi(page)
	newCount := paging(newPage, count)
	for _, val := range newCount {
		var goods models.GoodsItem
		db.Model(&goods).Where("id= ?", val).Find(&goods)
		newGoods := thecrawler.FormatConversion(goods)
		goodsinform = append(goodsinform, newGoods)
	}

	var item []interface{}
	item = append(item, goodsinform)
	item = append(item, map[string]int{"CounNumb": len(count)})
	_, err := ctx.JSON(item)

	if err != nil {
		log.Println("写入返回值失败：", err)
	}

}

//根据商品id进行搜索
func HandlerItemId(ctx iris.Context) {
	ctx.ResponseWriter().Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	ctx.ResponseWriter().Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	ctx.ResponseWriter().Header().Set("content-type", "application/json")             //返回数据格式是json

	itemId := ctx.FormValue("itemid")
	var goods models.GoodsItem
	db := models.Session
	db.Model(&goods).Where("down_type=?", 0).Where("Itemid= ?", itemId).Find(&goods)
	if goods.Itemid == "" {
		_, err := ctx.JSON(map[string]int{"errcode": 400})
		if err != nil {
			log.Println("写入返回值失败：")
		}
		return
	}
	newGoods := thecrawler.FormatConversion(goods)
	ctx.JSON(newGoods)
}

//分页处理
func paging(page int, count []int) []int {
	newAct := (page - 1) * 100
	end := page * 100
	if newAct > len(count) {
		page = len(count) - 1 - 100
		return nil
	}
	if end > len(count) {
		end = len(count) - 1
	}
	count = count[newAct:end]
	return count
}
