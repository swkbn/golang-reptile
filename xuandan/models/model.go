package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type GoodsItem struct {
	Id                       int       `gorm:"primary_key"`  //主键
	GoodsId                  int       `gotm:"unique_index"` //（唯一健）
	Itemid                   string    `gotm:"unique"`       //购买商品的ID	（唯一健）
	Itemtitle                string    //商品标题
	Itemshorttitle           string    //短标题
	Itemdesc                 string    //商品简介
	Itemprice                string    `gorm:"not null"` //商品原价
	Itemendprice             string    //券后价格
	Tkrates                  string    //佣金比例
	Tkmoney                  string    //预估佣金
	Itemsale                 string    //月销量
	Itemsale2                int       //近两小时销量
	Todaysale                string    //日销量
	Grade_avg                string    //宝贝评分
	Couponmoney              string    //优惠券价格
	Couponreceive            string    //已领取优惠券数量
	Couponnum                string    //优惠券共有多少张
	Couponstarttime          string    //优惠券开始时间
	Couponendtime            string    //优惠券结束时间
	Couponurl                string    //领券链接
	Activity_type            string    //活动类型
	Shoptype                 string    //购买类型
	Shopid                   string    //店铺id
	Shopname                 string    //店铺名称
	Copy_friends_circle_text []byte    `gorm:"type:varbinary(1000)"` //朋友圈推广文案
	Fqcat                    int       //商品的分类
	ItempicCopy              string    //商品展示的第二张图片（也是长图）
	Itempic                  string    //商品展示首长图
	Taobao_image_qiniu       string    //主图和推广图片
	Image                    string    `gorm:"type:varchar(500)"` //实拍图片
	Main_video_url           string    //视频链接
	Video_image              string    //视频图片
	Down_type                int       //判断商品是否下架
	Change_time              time.Time //更新时间
}

var Session *gorm.DB

//连接数据库，创建表
func init() {
	var err error
	//连接数据库
	Session, err = gorm.Open("mysql", "root:102030@/goods?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println("错误：", err)
	}

	//创建表
	if !Session.HasTable(&GoodsItem{}) {

		if err := Session.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&GoodsItem{}).Error; err != nil {
			log.Println(err)
		}
	}
	log.Println("连接数据库：。。。。。。。。。。。。。。。")
}
