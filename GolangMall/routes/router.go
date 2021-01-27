package routes

import (
	"cmall/api"
	"cmall/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

//路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("something-very-secret"))
	r.Use(middleware.Cors())
	r.Use(sessions.Sessions("mysession", store))
	v1 := r.Group("api/v1")
	{
		v1.POST("user/register", api.UserRegister) //用户注册
		v1.POST("user/login", api.UserLogin)       //用户登陆
		//邮箱绑定解绑接口
		v1.POST("user/vaild-email", api.VaildEmail)
		//商品操作
		v1.GET("products", api.ListProducts)
		v1.GET("products/:id", api.ShowProduct)
		v1.GET("carousels", api.ListCarousels)      //轮播图
		v1.GET("imgs/:id", api.ShowProductImgs)     //商品图片
		v1.GET("info-imgs/:id", api.ShowInfoImgs)   //商品详情图片操作
		v1.GET("param-imgs/:id", api.ShowParamImgs) //商品参数图片操作
		v1.GET("categories", api.ListCategories)    //商品分类操作
		v1.GET("rankings", api.ListRanking)         //排行榜
		v1.GET("elec-rankings", api.ListElecRanking)
		v1.GET("acce-rankings", api.ListAcceRanking)
		v1.POST("products", api.CreateProduct) //创建商品
		v1.GET("notices", api.ShowNotice)
		v1.POST("avatar", api.UploadToken) //上传操作
		authed := v1.Group("/")            //需要登陆保护
		authed.Use(middleware.JWT())
		{
			//authed.POST("products",api.CreateProduct)   //创建商品
			authed.GET("ping", api.CheckToken)               //验证token
			authed.PUT("user", api.UserUpdate)               //用户修改，更新
			authed.POST("user/sending-email", api.SendEmail) //用户操作
			//authed.POST("avatar",api.UploadToken) 		 //上传操作
			//收藏夹
			authed.GET("favorites/:id", api.ShowFavorites) //收藏夹操作
			authed.POST("favorites", api.CreateFavorite)
			authed.DELETE("favorites", api.DeleteFavorite)
			//购物车
			authed.POST("carts", api.CreateCart)
			authed.GET("carts/:id", api.ShowCarts)
			authed.PUT("carts", api.UpdateCart)
			authed.DELETE("carts", api.DeleteCart)
			//收获地址操作
			authed.POST("addresses", api.CreateAddress)
			authed.GET("addresses/:id", api.ShowAddresses)
			authed.PUT("addresses", api.UpdateAddress)
			authed.DELETE("addresses", api.DeleteAddress)
		}
	}
	v2 := r.Group("/api/v2")
	{
		//管理员登陆注册
		v2.POST("admin/register", api.AdminRegister)
		//登陆
		v2.POST("admin/login", api.AdminLogin)
		//商品操作
		v2.GET("products", api.ListProducts)
		v2.GET("products/:id", api.ShowProduct)
		//轮播图
		v2.GET("carousels", api.ListCarousels)
		//商品图片
		v2.GET("imgs/:id", api.ShowProductImgs)
		//分类操作
		v2.GET("categories", api.ListCarousels)

		v2.GET("users", api.ListUsers)
		authed2 := v2.Group("/")
		authed2.Use(middleware.JWTAdmin())
		{
			//商品操作
			authed2.POST("products", api.CreateProduct)
			//authed2.GET("users",api.ListUsers)
			authed2.DELETE("products/:id", api.DeleteProduct)
			authed2.PUT("products", api.UpdateProduct)
			//轮播图操作
			authed2.POST("carousels", api.CreateCarousel)
			//商品图片操作
			authed2.POST("imgs", api.CreateProductImg)
			//商品详情图片操作
			authed2.POST("info-imgs", api.CreateInfoImg)
			//商品参数图片操作
			authed2.POST("param-imgs", api.CreateParamImg)
			//分类操作
			authed2.POST("categories", api.CreateCategory)
			//公告操作
			authed2.POST("notices", api.CreateNotice)
			authed2.PUT("notice", api.UpdateNotice)
		}
	}
	return r
}
