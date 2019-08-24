package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro"

	"microproject/ihomeweb/utils"
	"context"
	getImageCode "microproject/getImageCode/proto/getImageCode"
	getArea "microproject/getArea/proto/getArea"
	getGetSmscd "microproject/getGetSmscd/proto/getGetSmscd"
	getpostRet "microproject/getpostRet/proto/getpostRet"
	getUserInfo"microproject/getUserInfo/proto/getUserInfo"
	getLogin "microproject/getLogin/proto/getLogin"
	getUploadphoto"microproject/getUploadphoto/proto/getUploadphoto"
	getChangeName "microproject/getChangeName/proto/getChangeName"
	getRealName "microproject/getRealName/proto/getRealName"
	getUserHouse "microproject/getUserHouse/proto/getUserHouse"
	getUserHousesInfo"microproject/getUserHousesInfo/proto/getUserHousesInfo"

	getUploadHouseImages"microproject/getUploadHouseImages/proto/getUploadHouseImages"

	getDetailedHouseMessage "microproject/getDetailedHouseMessage/proto/getDetailedHouseMessage"
	searchHouse "microproject/searchHouse/proto/searchHouse"
	"github.com/afocus/captcha"
	"encoding/json"
	"image/png"
	"regexp"

	"github.com/gin-contrib/sessions"

	"fmt"

	"path"

	"strconv"
)
//获取地域信息
func GetArea(ctx*gin.Context){
	//处理数据

	//consul
	reg:=consul.NewRegistry(func(options *registry.Options) {
		options.Addrs=[]string{"127.0.0.1:8500",}
	})

	//get service
	servicer:=grpc.NewService(
		micro.Registry(reg),
	)

	//get client
	client:=getArea.NewGetAreaService(utils.GetAreaName,servicer.Client())
	resp,err:=client.Call(context.TODO(),&getArea.Request{})

	//返回数据
	if err != nil {
		resp.Errno = utils.RECODE_DATAERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
	}

	ctx.JSON(200,resp)

}

//获取session
func GetSession(ctx*gin.Context){
	//存储session
	//设置当前ctx给浏览器传递cookie

	//回复数据容器

	resp:=make(map[string]interface{})
	temResp:=make(map[string]interface{})

	//获取session
	s:=sessions.Default(ctx)
	userName:=s.Get("userName")
	if userName==nil {
		//初始化返回值
		resp["errno"]=utils.RECODE_SESSIONERR
		resp["errmsg"]=utils.RecodeText(utils.RECODE_SESSIONERR)
	}else {
		resp["errno"]=utils.RECODE_OK
		resp["errmsg"]=utils.RecodeText(utils.RECODE_OK)
		temResp["name"]=userName.(string)
		resp["data"]=temResp
	}

	//返回值//200指的是http通信协议的无错误
	ctx.JSON(200,resp)
}

//获取图片验证码
func GetImageCd(ctx*gin.Context){
	//获取数据
	uuid:=ctx.Param("uuid")
	//校验数据//正则url代为处理
	//处理数据
	service:=utils.GetGrpcService()
	//获取客户端
	client:=getImageCode.NewGetImageCodeService(utils.GetImageName,service.Client())
	//调用远程
	resp,err:=client.Call(context.TODO(),&getImageCode.Request{Uuid:uuid})
	if err!=nil {
		resp.Errno=utils.RECODE_DATAERR
		resp.Errmsg=utils.RecodeText(utils.RECODE_DATAERR)
		ctx.JSON(200,resp)
		return
	}
	//获取数据,反序列化,获取到存入的image数据,然后传回前段
	var img captcha.Image

	json.Unmarshal(resp.Img,&img)

	png.Encode(ctx.Writer,img)

}

//获取验证码
func GetSmscd(ctx*gin.Context){
	//返回数据容器
	resp:=make(map[string]interface{})
	//获取数据

	phoneNum:=ctx.Param("mobile")
	text:=ctx.Query("text")
	id:=ctx.Query("id")
	//校验数据
	//校验电话号码格式
	reg,_:=regexp.Compile(`^1[3,4,5,7,8]\d{9}$`)
	result:=reg.MatchString(phoneNum)


	if !result {
		resp["errno"]=utils.RECODE_MOBILEERR
		resp["errmsg"]=utils.RecodeText(utils.RECODE_MOBILEERR)
		ctx.JSON(200,resp)
		return
	}

	//校验数据不能为空
	if text=="" || id=="" {
		resp["errno"] =utils.RECODE_DATAERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
		ctx.JSON(500,resp)
		return
	}
	//处理数据
	servicer:=utils.GetGrpcService()

	grpcClient:=getGetSmscd.NewGetGetSmscdService(utils.GetSesCd,servicer.Client())
	rsp,err:=grpcClient.Call( context.TODO(),&getGetSmscd.Request{
		Text:text,
		Uuid:id,
		PhoneNum:phoneNum,
	})
	//返回数据
	if err!=nil {
		resp["errno"]=utils.RECODE_SMSERR
		resp["errmsg"]=utils.RecodeText(utils.RECODE_SMSERR)
		ctx.JSON(200,resp)
		return
	}
	ctx.JSON(200,rsp)
}

type RegUser struct {
	Mobile string `json:"mobile"`
	PassWord string `json:"password"`
	SmsCode string `json:"sms_code"`
}
//注册服务
func PostRet(ctx*gin.Context){
	//获取数据
	//创建返回数据容器
	resp:=make(map[string]interface{})

	var user RegUser
	//由于前段规定的是json所以用bind
	//校验数据
	err:=ctx.Bind(&user)
	if err !=nil{
		resp["errno"] = utils.RECODE_REQERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_REQERR)
		ctx.JSON(200,resp)
		return
	}

	//处理数据
	//获取servicer
	servicer:=utils.GetGrpcService()

	grpcClient:=getpostRet.NewGetpostRetService(utils.GetpostRetName,servicer.Client())

	rsp,err:=grpcClient.Call(context.TODO(),&getpostRet.Request{
		PhoneNum:user.Mobile,
		Pwd:user.PassWord,
		SmsCode:user.SmsCode,
	})
	//返回数据
    fmt.Println(err)
	if err!=nil{
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
		ctx.JSON(200,resp)
		return
	}

	//注册成功添加session
	s:=sessions.Default(ctx)

	s.Set("userName",user.Mobile)
	err=s.Save()

	if err != nil {
		rsp.Errno = utils.RECODE_SESSIONERR

		rsp.Errmsg = utils.RecodeText(utils.RECODE_SESSIONERR)

	}

	ctx.JSON(200,rsp)
}

//退出业务
func DeleteSession(ctx*gin.Context){
	resp:=make(map[string]interface{})

	s:=sessions.Default(ctx)

	s.Delete("userName")
	err:=s.Save()
	if err!=nil {
		resp["errno"]=utils.RECODE_SESSIONERR
		resp["errmsg"]=utils.RecodeText(utils.RECODE_SESSIONERR)

	}else {
		resp["errno"]=utils.RECODE_OK
		resp["errmsg"]=utils.RecodeText(utils.RECODE_OK)
	}
	ctx.JSON(200,resp)
}

type UserLogin struct {
	Mobile string `json:"mobile"`
	PassWord string `json:"password"`
}
//登录
func PostLogin(ctx *gin.Context){
	//获取数据

	var userLogin UserLogin
	err:=ctx.Bind(&userLogin)
	//校验数据
	//错误容器
	errResp:=make(map[string]interface{})
	if err!=nil {
		errResp["errno"]=utils.RECODE_SESSIONERR
		errResp["errmsg"]=utils.RecodeText(utils.RECODE_SESSIONERR)
		ctx.JSON(200,errResp)
		return
	}

	//处理数据
	servicer:=utils.GetGrpcService()
	grpcClient:=getLogin.NewGetLoginService(utils.GetLogin,servicer.Client())
	resp,err:=grpcClient.Call(context.TODO(),&getLogin.Request{
		Mobile:userLogin.Mobile,
		Password:userLogin.PassWord,
	})
	if err != nil {
		fmt.Println(err)
		errResp["errno"] = utils.RECODE_LOGINERR
		errResp["errmsg"] = utils.RecodeText(utils.RECODE_LOGINERR)
		ctx.JSON(200,errResp)
		return
	}
	//返回数据
	if resp.Errno == utils.RECODE_OK {
		s:=sessions.Default(ctx)
		s.Set("userName",resp.Name)
		s.Save()
	}
	ctx.JSON(200,resp)
}

//获取登录信息
func GetUserInfo(ctx*gin.Context){
	//获取数据
	//获取登录信息,登录信息在session中存储
	s:=sessions.Default(ctx)

	//根据用户名获取用户信息
	userName:=s.Get("userName")

	//校验数据,并且创建错误容器
	errResp:=make(map[string]interface{})
	if userName==nil {
		errResp["errno"]=utils.RECODE_SESSIONERR
		errResp["errmsg"]=utils.RecodeText(utils.RECODE_SESSIONERR)
		ctx.JSON(200,errResp)
		return
	}

	//处理数据,并获取远端数据
	servicer:=utils.GetGrpcService()
	client5:=getUserInfo.NewGetUserInfoService(utils.GetUserInfo, servicer.Client())
	resp,err:=client5.Call(context.TODO(),&getUserInfo.Request{Name:userName.(string)})
	if err!=nil {
		fmt.Println(err)
		resp.Errno = utils.RECODE_DBERR
		resp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
	}

	//返回数据
	ctx.JSON(200,resp)
}

//上传头像业务
func PostAvatar(ctx*gin.Context){

	errResp:=make(map[string]interface{})
	//获取数据
	fileHeader,err:=ctx.FormFile("avatar")
	//校验数据
	if err != nil {
		errResp["errno"]=utils.RECODE_DATAERR
		errResp["errmsg"]=utils.RecodeText(utils.RECODE_DATAERR)
		ctx.JSON(200,errResp)
		return
	}

	//大小校验
	if fileHeader.Size > 50000{
		errResp["errno"]=utils.RECODE_FILEBIGER
		errResp["errmsg"]=utils.RecodeText(utils.RECODE_FILEBIGER)
		ctx.JSON(200,errResp)
		return
	}
	//格式校验
	fileExt:=path.Ext(fileHeader.Filename)
	if fileExt!=".jpg"&& fileExt!=".png"&&fileExt!=".jepg"{
		errResp["errno"]=utils.RECODE_FORMATERR
		errResp["errmsg"]=utils.RecodeText(utils.RECODE_FORMATERR)
		ctx.JSON(200,errResp)
		return
	}
	//获取文件留
	file,_:=fileHeader.Open()
	buffer:=make([]byte,fileHeader.Size)
	file.Read(buffer)
	//获取用户名
	s:=sessions.Default(ctx)
	userName:=s.Get("userName")
	//处理数据
	servicer:=utils.GetGrpcService()
	client4:=getUploadphoto.NewGetUploadphotoService(utils.GetUploadPhoto,servicer.Client())
	resp,err:=client4.Call(context.TODO(),&getUploadphoto.Request{
		FileExt:fileExt[1:],
		FileBuffer:buffer,
		Name:userName.(string),
	})

	if err != nil {
		fmt.Println(err)
		errResp["errno"] = utils.RECODE_DATAERR
		errResp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
		ctx.JSON(200,errResp)
		return
	}
	//返回数据
	ctx.JSON(200,resp)
}

//更新用户名
type userInfo struct {
	Name string `json:"name"`
}
//更改用户名
func PutUserInfo(ctx*gin.Context){
	//get data

	var user userInfo
	err:=ctx.Bind(&user)


	errResp:=make(map[string]interface{})
	if err!=nil {
		errResp["errno"]=utils.RECODE_REQERR
		errResp["errmsg"]=utils.RecodeText(utils.RECODE_REQERR)
		ctx.JSON(200,errResp)
		return
	}


	//从session中获取数据
	s:=sessions.Default(ctx)
	userName:=s.Get("userName")

	//处理数据
	servicer:=utils.GetGrpcService()
	client3:=getChangeName.NewGetChangeNameService(utils.GetChangePhoto,servicer.Client())
	resp,err:=client3.Call(context.TODO(),&getChangeName.Request{
		PreName:userName.(string),
		CurName:user.Name,
	})
	if err != nil {
		fmt.Println(err)
		errResp["errno"] = utils.RECODE_DATAERR
		errResp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
		ctx.JSON(200,errResp)
		return
	}
	if resp.Errno==utils.RECODE_OK {
		s.Set("userName",resp.Data.Name)
		s.Save()
	}

	ctx.JSON(200,resp)
	//back data
}

//展示用户实名认证
func GetUserAuth(ctx*gin.Context){}


type RealUser struct {
	RealName string `json:"real_name"`
	IdCard string `json:"id_card"`
}

func PutUserAuth(ctx*gin.Context){
	//获取数据
	var user RealUser
	err:=ctx.Bind(&user)
	errReap:=make(map[string]interface{})
	if err != nil {
		errReap["errno"]=utils.RECODE_REQERR
		errReap["errmsg"]=utils.RecodeText(utils.RECODE_REQERR)
		ctx.JSON(200,errReap)
		return
	}
	//校验数据
	req,err:=regexp.Compile(`^[1-9][0-7]\d{4}((19\d{2}(0[13-9]|1[012])(0[1-9]|[12]\d|30))|(19\d{2}(0[13578]|1[02])31)|(19\d{2}02(0[1-9]|1\d|2[0-8]))|(19([13579][26]|[2468][048]|0[48])0229))\d{3}(\d|X|x)?$`)
	if err!=nil{
		fmt.Println(err)
		return
	}
	if !req.MatchString(user.IdCard) {

		errReap["errno"]=utils.RECODE_IDCARERR
		errReap["errmsg"]=utils.RecodeText(utils.RECODE_IDCARERR)
		ctx.JSON(200,errReap)
		return
	}
	//校验身份证信息是否真实,调用外部接口实现身份证校验,推荐使用聚合

	//从session中获取用户名
	s:=sessions.Default(ctx)
	userName:=s.Get("userName")
	//处理数据
	//服务调用
	servicer:=utils.GetGrpcService()
	client2:=getRealName.NewGetRealNameService(utils.GetRealName,servicer.Client())

	//调用远程服务
	resp,err:=client2.Call(context.TODO(),&getRealName.Request{
		Name:userName.(string),
		RealName:user.RealName,
		IdCard:user.IdCard,
	})

	//返回数据
	if err != nil {
		fmt.Println(err)
		resp.Errno=utils.RECODE_SERVERERR
		resp.Errmsg=utils.RecodeText(utils.RECODE_SERVERERR)
		ctx.JSON(200,resp)
		return
	}
	ctx.JSON(200,resp)
}

//获取已发布房源信息
func GetUserHouses(ctx*gin.Context){
	//获取数据
	s:=sessions.Default(ctx)
	userName:=s.Get("userName")
	//调用远程服务
	grpcService := utils.GetGrpcService()
	serclient:=getUserHousesInfo.NewGetUserHousesInfoService(utils.GetUserHousesInfo,grpcService.Client())

	resp,err := serclient.Call(context.TODO(),&getUserHousesInfo.Request{Name:userName.(string)})


	if err != nil {
		fmt.Println(err)
		resp.Errno=utils.RECODE_SERVERERR
		resp.Errmsg=utils.RecodeText(utils.RECODE_SERVERERR)
	}

	ctx.JSON(200,resp)
}

type UserMessage struct {
	Acreage  string 	 	`json:"acreage"`
	Address string 			`json:"address"`
	AreaId string				`json:"area_id"`
	Beds string				`json:"beds"`
	Capacity string			`json:"capacity"`
	Deposit string			`json:"deposit"`
	MaxDays string				`json:"max_days"`
	MinDays string			`json:"min_days"`
	Price string				`json:"price"`
	RoomCount string		`json:"room_count"`
	Title string			`json:"title"`
	Unit string		 		`json:"unit"`
	Facility []string		`json:"facility"`
}
//发布房源
func PostHouses(ctx*gin.Context){
	//获取数据
	var userKK UserMessage
	err:=ctx.Bind(&userKK)
	//校验数据
	errReap:=make(map[string]interface{})
	if err != nil {
		errReap["errno"]=utils.RECODE_REQERR
		errReap["errmsg"]=utils.RecodeText(utils.RECODE_REQERR)
		ctx.JSON(200,errReap)
		return
	}

	//从session中获取用户名
	s:=sessions.Default(ctx)
	userName:=s.Get("userName")
	//处理数据
	servicer:=utils.GetGrpcService()
	client1:=getUserHouse.NewGetUserHouseService(utils.GetUserHouse,servicer.Client())
	resp,err:=client1.Call(context.TODO(),&getUserHouse.Request{
		Name:userName.(string),
		Acreage:userKK.Acreage,
		Address:userKK.Address,
		AreaId :userKK.AreaId,
		Beds :userKK.Beds,
		Capacity :userKK.Capacity,
		Deposit  :userKK.Deposit,
		MaxDays  :userKK.MaxDays,
		MinDays  :userKK.MinDays,
		Price    :userKK.Price,
		RoomCount :userKK.RoomCount,
		Title  :userKK.Title,
		Unit :userKK.Unit,
		Facility :userKK.Facility,
	})
	//返回数据
	if err != nil {
		fmt.Println(err)
		resp.Errno=utils.RECODE_SERVERERR
		resp.Errmsg=utils.RecodeText(utils.RECODE_SERVERERR)
		ctx.JSON(200,resp)
		return
	}
	ctx.JSON(200,resp)

}

//上传房源图片
func PostHousesImage(ctx*gin.Context){
	//获取数据
	//1.获取url传过来的Id
	Id:=ctx.Param("id")
	//2.获取图片和文件的格式
	errResp:=make(map[string]interface{})
	//校验数据
	filebuffer,fileExt:=utils.UploadPicture(ctx,errResp,"house_image")
	if filebuffer==nil {
		ctx.JSON(200,errResp)
		return
	}
	//处理数据
	id,_:=strconv.Atoi(Id)
	//1.调用远程服务
	servicer:=utils.GetGrpcService()

	client:=getUploadHouseImages.NewGetUploadHouseImagesService(utils.GetUploadHouseImage,servicer.Client())
	resp,err:=client.Call(context.TODO(),&getUploadHouseImages.Request{
		HouseId:int32(id),
		HouseBuffer:filebuffer,
		FileExt:fileExt,
	})
	//返回数据
	if err != nil {
		fmt.Println(err)
		errResp["errno"] = utils.RECODE_SERVERERR
		errResp["errmsg"] = utils.RecodeText(utils.RECODE_SERVERERR)
		ctx.JSON(200,errResp)
		return
	}
	ctx.JSON(200,resp)
}

//获取房源详细信息
func GetHouseInfo(ctx*gin.Context){
	//获取数据
	houseId:=ctx.Param("id")

	errResp:=make(map[string]interface{})

	//返回数据
	if houseId=="" {

		errResp["errno"]=utils.RECODE_REQERR
		errResp["errmsg"]=utils.RecodeText(utils.RECODE_REQERR)
		ctx.JSON(200,errResp)
		return
	}
	//处理数据
	servicer:=utils.GetGrpcService()
	client:=getDetailedHouseMessage.NewGetDetailedHouseMessageService(utils.GetDetailHouse,servicer.Client())
	resp,err:=client.Call(context.TODO(),&getDetailedHouseMessage.Request{
		HouseId:houseId,
	})
	fmt.Println(resp)
	//返回数据
	if err != nil {
		fmt.Println(err)
		errResp["errno"]=utils.RECODE_SERVERERR
		errResp["errmsg"]=utils.RecodeText(utils.RECODE_SERVERERR)
		ctx.JSON(200,errResp)
		return
	}

	ctx.JSON(200,resp)
}

//搜索房屋
func GetHouses(ctx* gin.Context){
	errResp := make(map[string]interface{})
	//获取数据
	aid := ctx.Query("aid")
	startTime := ctx.Query("sd")
	endTime := ctx.Query("ed")
	//数据校验
	if aid == "" || startTime == "" || endTime ==""{
		errResp["errno"] = utils.RECODE_REQERR
		errResp["errmsg"] = utils.RecodeText(utils.RECODE_REQERR)
		ctx.JSON(200,errResp)
		return
	}

	//处理数据
	grpcServicer:=utils.GetGrpcService()

	client:=searchHouse.NewSearchHouseService(utils.GetSearchHouse,grpcServicer.Client())

	resp,err:=client.Call(context.TODO(),&searchHouse.Request{
		Aid:aid,
		StartTime:startTime,
		EndTime:endTime,
		CityAreaOptions:"",
	})
	if err != nil {
		fmt.Println(err)
		errResp["errno"] = utils.RECODE_SERVERERR
		errResp["errmsg"] = utils.RecodeText(utils.RECODE_SERVERERR)
		ctx.JSON(200,errResp)
		return
	}
	ctx.JSON(200,resp)


}
