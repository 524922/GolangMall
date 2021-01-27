package model

import (
	"cmall/cache"
	"github.com/jinzhu/gorm"
	"strconv"
)

//商品模型
type Product struct {
	gorm.Model
	ProductID     string `gorm:"primary_key"`
	Name          string
	CategoryID    int
	Title         string
	Info          string `gorm:"size:1000"`
	ImgPath       string
	Price         string
	DiscountPrice string
	BossID        int
	BossName      string
	BossAvatar    string
}

func (Product *Product) View() uint64 {
	//增加视频点击数
	countStr, _ := cache.RedisClient.Get(cache.ProductViewKey(Product.ID)).Result()
	count, _ := strconv.ParseUint(countStr, 10, 64)
	return count
}

//AddView 视频游览
func (Product *Product) AddView() {
	//增加视频点击数
	cache.RedisClient.Incr(cache.ProductViewKey(Product.ID))
	//增加排行点击数
	cache.RedisClient.ZIncrBy(cache.RankKey, 1, strconv.Itoa(int(Product.ID)))
}

//AddElecRank 增加加点排行点击数
func (Product *Product) AddElecRank() {
	//增加家电排汗点击数
	cache.RedisClient.ZIncrBy(cache.ElectricalRank, 1, strconv.Itoa(int(Product.ID)))
}

//AddAcceRank 增加配件排行点击数
func (Product *Product) AddAcceRank() {
	//增加配件排行点击数
	cache.RedisClient.ZIncrBy(cache.AccessoryRank, 1, strconv.Itoa(int(Product.ID)))
}
