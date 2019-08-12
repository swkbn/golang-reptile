package thecrawler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"xuandan/models"
)

//创建用于接收数据的结构体

type Preview struct {
	Data Data `json:"data"`
}
type Data struct {
	Front              []Front      `json:"front"` //请求首页中的数据
	Back               []Back       `json:"back"`
	Fqcat              string       `json:"fqcat"`              //商品分类
	Itemid             string       `json:"itemid"`             //购买商品的ID
	Itemtitle          string       `json:"itemtitle"`          //商品标题
	Itemshorttitle     string       `json:"itemshorttitle"`     //短标题
	Itemdesc           string       `json:"itemdesc"`           //商品简介
	ItempicCopy        string       `json:"itempic_copy"`       //商品展示的第二张图片（也是长图）
	Itempic            string       `json:"itempic"`            //商品展示首长图
	Itemprice          string       `json:"itemprice"`          //商品原价
	Itemendprice       string       `json:"itemendprice"`       //券后价格
	Tkrates            string       `json:"tkrates"`            //佣金比例
	Tkmoney            string       `json:"tkmoney"`            //预估佣金
	Itemsale           string       `json:"itemsale"`           //月销量
	Itemsale2          string       `json:"itemsale2"`          //近两小时销量
	Todaysale          string       `json:"todaysale"`          //日销量
	Grade_avg          string       `json:"grade_avg"`          //宝贝评分
	Couponmoney        string       `json:"couponmoney"`        //优惠券价格
	Couponreceive      string       `json:"couponreceive"`      //已领取优惠券数量
	Couponnum          string       `json:"couponnum"`          //优惠券共有多少张
	Shopid             string       `json:"shopid"`             //店铺id
	Shopname           string       `json:"shopname"`           //店铺名称
	Shoptype           string       `json:"shoptype"`           //	购买类型
	Couponstarttime    string       `json:"couponstarttime"`    //优惠券开始时间
	Couponendtime      string       `json:"couponendtime"`      //优惠券结束时间
	Couponurl          string       `json:"couponurl"`          //领券链接
	Activity_type      string       `json:"activity_type"`      //活动类型
	Videoid            string       `json:"videoid"`            //视频id
	Taobao_image_qiniu string       `json:"taobao_image_qiniu"` //主图和推广图片
	Down_type          string       `json:"down_type"`          //判断商品是否下架
	MaterialInfo       MaterialInfo `json:"material_info"`      //下一个结构体
}
type Front struct {
	Id    string `json:"id"`
	Fqcat string `json:"fqcat"`
}
type Back struct {
	Id    string `json:"id"`
	Fqcat string `json:"fqcat"`
}
type MaterialInfo struct {
	Image                    string `json:"image"`                    //实拍图片
	Main_video_url           string `json:"main_video_url"`           //视频链接
	VideoImage               string `json:"video_image"`              //视频图片
	Friends_circle_text      string `json:"friends_circle_text"`      //朋友圈文案
	Copy_friends_circle_text string `json:"copy_friends_circle_text"` //朋友圈文案
	Show_friends_circle_text string `json:"show_friends_circle_text"` //朋友圈文案
	Couponlife               string `json:"couponlife"`               //有效期
}

var m sync.RWMutex
var ms sync.Mutex

//错误捕获
func recoverName() {
	if r := recover(); r != nil {
		log.Println("recovered from ", r)
	}
}

//部分页数
func PartPage(act, end int) {
	defer recoverName()
	for i := act; i <= end; i++ {
		resp, err := http.Get("https://www.haodanku.com/indexapi/hdk_list?type=1&p=" + strconv.Itoa(i) + "&search_type=0&category_id=0&price_min=&price_max=&array_type=&sale_min=&tkrates_min=&coupon_max=&tkmoney_min=&avg_min=&discount_max=")
		if err != nil {
			log.Println("请求错误：", err)
		}
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("读取错误", err)
		}
		var prev Preview
		err = json.Unmarshal(buf, &prev)
		if err != nil {
			log.Println("转换json错误", err)
		}
		go disposePrev(i, prev)
	}
	myTicker := time.NewTicker(time.Minute * 5)
	<-myTicker.C
	go PartPage(act, end)
	log.Println("更新前10页内容")
	runtime.Goexit()
}

//全部页码
func AllPage() {
	defer recoverName()
	for i := 1; ; i++ {
		resp, err := http.Get("https://www.haodanku.com/indexapi/hdk_list?type=1&p=" + strconv.Itoa(i) + "&search_type=0&category_id=0&price_min=&price_max=&array_type=&sale_min=&tkrates_min=&coupon_max=&tkmoney_min=&avg_min=&discount_max=")
		if err != nil {
			log.Println("请求错误：", err)
		}
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("读取错误", err)
		}
		var prev Preview
		err = json.Unmarshal(buf, &prev)
		if err != nil {
			log.Println("转换json错误", err)
		}
		back := prev.Data.Back
		if back == nil {
			return
		}
		go disposePrev(i, prev)
	}
	myTicker := time.NewTicker(time.Hour * 6)
	<-myTicker.C
	log.Println("请求全部数据")
	go AllPage()
}

//数据处理
func disposePrev(page int, prev Preview) {
	defer recoverName()
	if page == 1 {
		for _, val := range prev.Data.Front {
			getDate(val.Id, val.Fqcat)
		}
	}
	for _, val := range prev.Data.Back {
		getDate(val.Id, val.Fqcat)
	}
	runtime.Goexit()
}

//获取数据(单个商品请求的话直接调用这个函数就可以)
func getDate(id, fqcat string) {
	defer recoverName()
	resp, err := http.Get("https://www.haodanku.com/detail/item_info?id=" + id)
	if err != nil {
		log.Println("请求错误", err)
	}
	buf, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("读取错误", err)
	}
	var data Preview
	json.Unmarshal(buf, &data)

	//请求视频链接，并赋值给结构体
	if data.Data.Videoid == "0" {
		url := getVideoUrl(data.Data.Itemid)
		data.Data.MaterialInfo.Main_video_url = url
	} else {
		data.Data.MaterialInfo.Main_video_url = "http://cloud.video.taobao.com/play/u/1/p/1/e/6/t/1/" + data.Data.Videoid + ".mp4"
	}
	VideoImage := strings.Split(data.Data.Taobao_image_qiniu, ",")
	if data.Data.Taobao_image_qiniu != "" {
		data.Data.MaterialInfo.VideoImage = "http://img.haodanku.com/" + VideoImage[0]
	}
	data.Data.Fqcat = fqcat
	dataStorage(data.Data, id)
}

//把数据存放到数据库中
func dataStorage(data Data, id string) {
	defer recoverName()
	db := models.Session
	var goodsInform models.GoodsItem
	var count int
	goodsId, _ := strconv.Atoi(id)
	goodsInform.GoodsId = goodsId
	fqcat, _ := strconv.Atoi(data.Fqcat)
	goodsInform.Fqcat = fqcat
	goodsInform.Shopid = data.Shopid
	goodsInform.Shopname = data.Shopname
	goodsInform.Itemid = data.Itemid
	goodsInform.Itemtitle = data.Itemtitle
	goodsInform.Itemshorttitle = data.Itemshorttitle
	goodsInform.Itemdesc = data.Itemdesc
	goodsInform.Itemprice = data.Itemprice
	goodsInform.Itemendprice = data.Itemendprice
	goodsInform.Tkrates = data.Tkrates
	goodsInform.Tkmoney = data.Tkmoney
	goodsInform.Itemsale = data.Itemsale
	Itemsale2, _ := strconv.Atoi(data.Itemsale2)
	goodsInform.Itemsale2 = Itemsale2
	goodsInform.Todaysale = data.Todaysale
	goodsInform.Grade_avg = data.Grade_avg
	goodsInform.Couponmoney = data.Couponmoney
	goodsInform.Couponreceive = data.Couponreceive
	goodsInform.Couponnum = data.Couponnum
	goodsInform.Couponstarttime = data.Couponstarttime
	goodsInform.Couponendtime = data.Couponendtime
	goodsInform.Couponurl = data.Couponurl
	goodsInform.Activity_type = data.Activity_type
	goodsInform.Shoptype = data.Shoptype
	//朋友圈文案
	goodsInform.Copy_friends_circle_text = []byte(filter(data.MaterialInfo.Copy_friends_circle_text))
	goodsInform.ItempicCopy = data.ItempicCopy
	goodsInform.Itempic = data.Itempic
	goodsInform.Taobao_image_qiniu = data.Taobao_image_qiniu
	goodsInform.Image = data.MaterialInfo.Image
	goodsInform.Main_video_url = data.MaterialInfo.Main_video_url
	goodsInform.Video_image = data.MaterialInfo.VideoImage
	down_type, _ := strconv.Atoi(data.Down_type)
	goodsInform.Down_type = down_type
	goodsInform.Change_time = time.Now()
	//先查询是否存在 然后在插入 或者更新
	ms.Lock()
	db.Model(&models.GoodsItem{}).Where("goods_id = ?", id).Count(&count)

	if count == 0 {
		db.Model(&models.GoodsItem{}).Create(&goodsInform)
	} else {
		//更新数据
		db.Model(&models.GoodsItem{}).Where("goods_id = ?", id).Updates(&goodsInform)
	}
	ms.Unlock()
}

//过滤字符串中的特殊字符
func filter(str string) string {
	bo := strings.Contains(str, "img")
	str = strings.Replace(str, `<br>`, "", -1)
	var copy_str string
	if bo {
		str = strings.Replace(str, "<", ">", -1)
		aa := strings.Split(str, ">")
		for _, val := range aa {
			bo := strings.Contains(val, "img")
			if !bo {
				copy_str = copy_str + val
			}
		}
		return copy_str
	} else {
		return str
	}
}

//请求视频链接
func getVideoUrl(itemid string) string {
	defer recoverName()
	resp, _ := http.Get("https://detail.tmall.com/item.htm?id=" + itemid)

	if resp == nil {
		return ""
	}
	buf, _ := ioutil.ReadAll(resp.Body)

	reg := regexp.MustCompile(`"imgVedioUrl":"(.*?).swf`)
	buffer := reg.FindAllStringSubmatch(string(buf), -1)
	if buffer != nil {
		url := buffer[0][1] + ".mp4"
		url = strings.Replace(url, "e/1", "e/6", -1)
		return url
	}
	return ""
}

//用于刷新已经过期的商品
func OverdueGoods() {
	defer recoverName()
	db := models.Session
	var count []int
	db.Model(&models.GoodsItem{}).Order("itemsale2 desc", true).Where("down_type = ?", 0).Pluck("goods_id", &count)
	go updateOverdue(count)
	for _, val := range count {
		updateOverdueGoods(strconv.Itoa(val))
		//用来控制协程的数量
	}
	go OverdueGoods()
	runtime.Goexit()
}

//只更新前300个
func updateOverdue(count []int) {
	if len(count) <= 300 {
		return
	}
	counter := count[:300]
	for _, val := range counter {
		updateOverdueGoods(strconv.Itoa(val))
	}
	tick := time.NewTicker(time.Minute * 5)
	<-tick.C
	updateOverdue(count)

}

//更新商品过期字段
func updateOverdueGoods(id string) {
	defer recoverName()
	var goodsInform Preview
	var db = models.Session
	resp, err := http.Get("https://www.haodanku.com/detail/item_info?id=" + id)
	if err != nil {
		log.Println("请求错误：", err)
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(buf, &goodsInform)
	if err != nil {
		m.Lock()
		db.Delete(models.GoodsItem{}, "goods_id =?", id)
		m.Unlock()

	}
	down_type, _ := strconv.Atoi(goodsInform.Data.Down_type)
	goodsid, _ := strconv.Atoi(id)
	m.Lock()
	err = db.Model(&models.GoodsItem{}).Where("goods_id = ?", goodsid).Update("down_type", down_type).Error
	m.Unlock()
	if err != nil {
		log.Println("出错在哪里：", err)
	}

	runtime.Goexit()
}
