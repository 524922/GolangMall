package serializer

import "cmall/model"

//购物车
type Cart struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"user_id"`
	ProductID     uint   `json:"product_id"`
	CreateAt      int64  `json:"create_at"`
	Num           uint   `json:"num"`
	MaxNum        uint   `json:"max_num"`
	Check         bool   `json:"check"`
	Name          string `json:"name"`
	ImgPath       string `json:"img_path"`
	DiscountPrice string `json:"discount_price"`
	Boss          string `json:"boss"`
}

func BuildCart(item1 model.Cart, item2 model.Product, item3 model.Boss) Cart {
	return Cart{
		ID:            item1.ID,
		UserID:        item1.UserID,
		ProductID:     item1.ProductID,
		CreateAt:      item1.CreatedAt.Unix(),
		Num:           item1.Num,
		MaxNum:        item1.MaxNum,
		Check:         false,
		Name:          item2.Name,
		ImgPath:       item2.ImgPath,
		DiscountPrice: item2.DiscountPrice,
		Boss:          item3.UserName,
	}
}

func BuildCarts(items []model.Cart) (carts []Cart) {
	for _, item1 := range items {
		item2 := model.Product{}
		item3 := model.Boss{}
		err := model.DB.First(&item2, item1.ProductID).Error
		if err != nil {
			continue
		}
		cart := BuildCart(item1, item2, item3)
		carts = append(carts, cart)
	}
	return carts
}
