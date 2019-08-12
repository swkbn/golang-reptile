package thecrawler

import (
	"strconv"
	"strings"
	"time"
	"xuandan/models"
)

type Goods struct {
	Id                  string     `json:"id"`                  //商品ID
	Goodsid             int        `json:"goodsid"`             //id
	Title               string     `json:"title"`               //标题
	Short_title         string     `json:"short_title"`         //短标题
	Image               string     `json:"image"`               //图片(商品展示第一张图片)
	Images              string     `json:"images"`              //长图片图片（也是商品展示第二张图片）
	Price               float64    `json:"price"`               //价格
	Origin_price        float64    `json:"origin_price"`        //原价
	Text                string     `json:"text"`                //简介
	Commission_rate     float64    `json:"commission_rate"`     //佣金比例
	Commission_money    float64    `json:"commission_money"`    //佣金价格
	Activity_type       string     `json:"activity_type"`       //活动类型
	Activity_id         string     `json:"activity_id"`         //领券ID
	Coupon_price        int        `json:"coupon_price"`        //优惠券价格
	Coupon_remain_count int        `json:"coupon_remain_count"` //优惠券剩余数量
	Coupon_started_at   time.Time  `json:"coupon_started_at"`   //优惠券开始时间
	Coupon_finished_at  time.Time  `json:"coupon_finished_at"`  //优惠券结束时间
	Coupon_total        int        `json:"coupon_total"`        //优惠券总数量
	Coupon_picked       int        `json:"coupon_picked"`       //优惠券已领取量
	Sales_count         int        `json:"sales_count"`         //销量
	Sales_2hours_count  int        `json:"sales_2_hours_count"` //两小时销量
	Todaysale           int        `json:"todaysale"`           //日销量
	Shop_type           string     `json:"shop_type"`           //店铺类型
	Shop_id             int        `json:"shop_id"`             //店铺ID
	Shop_name           string     `json:"shop_name"`           //店铺名称
	Provider            string     `json:"provider"`            //供应商？？？
	Sns_content         string     `json:"sns_content"`         //推广内容
	Sns_image           []string   `json:"sns_image"`           //推广图片集
	Contents            []Contents `json:"contents"`
	Updated_time        time.Time  `json:"updated_time"` //更新时间
}
type Contents struct {
	Type  string `json:"type"` //（带有image类型的是实拍图）（还有video是实拍视频）
	Value string `json:"value"`
	Thumb string `json:"thumb"`
}

//返回数据的格式转换
func FormatConversion(gsInfo models.GoodsItem) Goods {
	var goods Goods
	goods.Goodsid = gsInfo.GoodsId
	goods.Id = gsInfo.Itemid
	goods.Title = gsInfo.Itemtitle
	goods.Short_title = gsInfo.Itemshorttitle
	goods.Image = gsInfo.Itempic
	//推广图片
	images := strings.Split(gsInfo.Taobao_image_qiniu, ",")
	for _, val := range images {
		goods.Sns_image = append(goods.Sns_image, "http://img.haodanku.com/"+val)
	}
	Itemendprice, _ := strconv.ParseFloat(gsInfo.Itemendprice, 64)
	goods.Price = Itemendprice
	Itemprice, _ := strconv.ParseFloat(gsInfo.Itemprice, 64)
	goods.Origin_price = Itemprice
	goods.Text = gsInfo.Itemdesc
	Tkrates, _ := strconv.ParseFloat(gsInfo.Tkrates, 64)
	goods.Commission_rate = Tkrates
	Tkmoney, _ := strconv.ParseFloat(gsInfo.Tkmoney, 64)
	goods.Commission_money = Tkmoney
	goods.Activity_type = gsInfo.Activity_type
	buf := strings.Split(gsInfo.Couponurl, "=")
	goods.Activity_id = buf[2]
	Couponmoney, _ := strconv.Atoi(gsInfo.Couponmoney)
	goods.Coupon_price = Couponmoney
	Couponstarttime, _ := strconv.ParseInt(gsInfo.Couponstarttime, 10, 64)
	goods.Coupon_started_at = time.Unix(Couponstarttime, 0)
	couponendtime, _ := strconv.ParseInt(gsInfo.Couponendtime, 10, 64)
	goods.Coupon_finished_at = time.Unix(couponendtime, 0)
	Couponnum, _ := strconv.Atoi(gsInfo.Couponnum)
	goods.Coupon_total = Couponnum
	Couponreceive, _ := strconv.Atoi(gsInfo.Couponreceive)
	goods.Coupon_remain_count = goods.Coupon_total - Couponreceive //优惠券剩余数量
	goods.Coupon_picked = Couponreceive                            //已领取
	Itemsale, _ := strconv.Atoi(gsInfo.Itemsale)
	goods.Sales_count = Itemsale
	goods.Sales_2hours_count = gsInfo.Itemsale2
	goods.Shop_type = gsInfo.Shoptype
	Shopid, _ := strconv.Atoi(gsInfo.Shopid)
	today, _ := strconv.Atoi(gsInfo.Todaysale)
	goods.Todaysale = today
	goods.Shop_id = Shopid //店铺ID
	goods.Shop_name = gsInfo.Shopname
	goods.Provider = "好单库" //供应商？？
	//朋友圈文案
	copy_frirnds := string(gsInfo.Copy_friends_circle_text)
	//copy_frirnds=strings.Replace(copy_frirnds,`<br>`,"",-1)
	goods.Sns_content = copy_frirnds

	goods.Images = "http://img.haodanku.com/" + gsInfo.ItempicCopy //长图片
	var conts Contents

	image := strings.Split(gsInfo.Image, ",")
	for _, val := range image {
		conts = Contents{"image", val, ""}
		goods.Contents = append(goods.Contents, conts)
	}
	conts = Contents{"video_url", gsInfo.Main_video_url, gsInfo.Video_image}
	goods.Contents = append(goods.Contents, conts)

	goods.Updated_time = gsInfo.Change_time

	return goods
}
