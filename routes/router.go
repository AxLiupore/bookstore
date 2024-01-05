package routes

import (
	api "bookstore/api/v1"
	"bookstore/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Cors())
	r.StaticFS("/static", http.Dir("./static"))

	v1 := r.Group("api/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "success")
		})

		// 用户操作
		v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.UserLogin)

		authed := v1.Group("/") // 需要登录保护
		authed.Use(middleware.JWT())
		{
			// 用户操作
			authed.POST("user/update", api.UserUpdate)
			authed.POST("user/avatar", api.UploadAvatar)
			authed.POST("user/send-email", api.SendEmail)
		}
	}
	return r
}
