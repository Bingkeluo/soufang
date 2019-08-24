package utils

import (
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/service/grpc"
	"github.com/micro/go-micro"
	"github.com/gin-gonic/gin"

	"path"
)

func GetGrpcService()micro.Service {

	reg:=consul.NewRegistry(func(options *registry.Options) {
		options.Addrs=[]string{
			"127.0.0.1:0:8500",
		}
	})

	servicer:=grpc.NewService(
		micro.Registry(reg),
	)

	return servicer
}

//上传图片
func UploadPicture(ctx *gin.Context,errResp map[string]interface{},fileImg string)([]byte,string){

	//获取数据
	fileHeader,err:=ctx.FormFile(fileImg)
	//校验数据
	if err != nil {
		errResp["errno"]=RECODE_DATAERR
		errResp["errmsg"]=RecodeText(RECODE_DATAERR)

		return nil,""
	}

	//大小校验
	if fileHeader.Size > 50000{
		errResp["errno"]=RECODE_FILEBIGER
		errResp["errmsg"]=RecodeText(RECODE_FILEBIGER)
		return nil,""
	}
	//格式校验
	fileExt:=path.Ext(fileHeader.Filename)
	if fileExt!=".jpg"&& fileExt!=".png"&&fileExt!=".jepg"{
		errResp["errno"]=RECODE_FORMATERR
		errResp["errmsg"]=RecodeText(RECODE_FORMATERR)

		return nil,""
	}
	//获取文件流
	file,_:=fileHeader.Open()
	buffer:=make([]byte,fileHeader.Size)
	file.Read(buffer)
	return buffer,fileExt

}