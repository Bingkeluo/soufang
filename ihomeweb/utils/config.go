package utils


//redis ip  port
const Redis_Address  = "192.168.21.81"
const Redis_Port  = "6379"

//mysql  ip   port
const MySQL_Address  = "127.0.0.1"
const MySQL_Port  = "3306"
const MySQL_UserName  = "root"
const MySQL_Pwd = "123456"
const MySQL_DB = "searchHouse"

//运行端口


//服务名
const GetAreaName = "go.micro.srv.getArea"
const GetImageName  = "go.micro.srv.getImageCode"
const GetSesCd  = "go.micro.srv.getGetSmscd"
const GetpostRetName  = "go.micro.srv.getpostRet"
const GetUserInfo  ="go.micro.srv.getUserInfo"
const GetLogin  = "go.micro.srv.getLogin"
const GetUploadPhoto  ="go.micro.srv.getUploadphoto"
const GetChangePhoto  = "go.micro.srv.getChangeName"
const GetRealName  = "go.micro.srv.getRealName"
const GetUserHouse  = "go.micro.srv.getUserHouse"
const GetUserHousesInfo  = "go.micro.srv.getUserHousesInfo"
const GetUploadHouseImage  = "go.micro.srv.getUploadHouseImages"
const GetDetailHouse  = "go.micro.srv.getDetailedHouseMessage"
const GetSearchHouse  = "go.micro.srv.searchHouse"
//上传图片所用的url
const NginxPath  ="http://192.168.21.81:8888/"