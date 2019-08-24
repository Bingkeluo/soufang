package main

import (
	"github.com/gin-gonic/gin"
	"microproject/ihomeweb/model"
	"microproject/ihomeweb/controller"
	"github.com/gin-contrib/sessions/redis"
	"microproject/ihomeweb/utils"
	"github.com/gin-contrib/sessions"
)


//实现中间件
func IsLogin()gin.HandlerFunc{
	return func(ctx *gin.Context) {
		s:=sessions.Default(ctx)
		userName:=s.Get("userName")
		if userName==nil {
			resp:=make(map[string]interface{})
			resp["errno"]=utils.RECODE_SESSIONERR
			resp["errmsg"]=utils.RecodeText(utils.RECODE_SESSIONERR)
			ctx.JSON(200,resp)
		}else{
			ctx.Next()
		}
	}

}
func main(){
	//获取路有对象
	router := gin.Default()

	//初始化
	//使用redis做session的存儲
	store,_:=redis.NewStore(10,"tcp",utils.Redis_Address+":"+utils.Redis_Port,"",[]byte("ihome"))

	//設置臨時的session
	store.Options(sessions.Options{
		MaxAge:60*60*24,
	})
	//使用session
	router.Use(sessions.Sessions("mysession",store))

	//路由匹配
	router.Static("/home","view")


	model.InitDb()

	//REST概念   路由风格
	r1 := router.Group("/api/v1.0")
	{
		r1.GET("/areas",controller.GetArea)
		r1.GET("/session",controller.GetSession)
		r1.GET("/imagecode/:uuid",controller.GetImageCd)
		r1.GET("/smscode/:mobile",controller.GetSmscd)
		r1.POST("/users",controller.PostRet)
		r1.DELETE("/session",controller.DeleteSession)

		r1.POST("/sessions",controller.PostLogin)

		r1.Use(IsLogin())
		r1.GET("/user",controller.GetUserInfo)

		//上传用户头像
		r1.POST("/user/avatar",controller.PostAvatar)
		r1.PUT("/user/name",controller.PutUserInfo)
		//获取用户实名信息
		r1.GET("/user/auth",controller.GetUserInfo)
		//更新实名认证
		r1.POST("/user/auth",controller.PutUserAuth)
		//获取房源信息
		r1.GET("/user/houses",controller.GetUserHouses)
		//发布房源信息
		r1.POST("/houses",controller.PostHouses)
		//上传房源图片
		r1.POST("/houses/:id/images",controller.PostHousesImage)
		//获取房屋详细信息
		r1.GET("/houses/:id",controller.GetHouseInfo)
		//搜索房屋
		r1.GET("/houses",controller.GetHouses)


	}

	//开启服务
	router.Run(":8081")
}
