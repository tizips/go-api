package helper

import (
	"bufio"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/gookit/goutil/dump"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
	"os"
	"path"
	"saas/kernel/config"
	"saas/kernel/config/configs"
	"strings"
)

func DoUploadBySimple(ctx *gin.Context, dirs string, file *multipart.FileHeader) (error, *UploadBySimple) {

	if config.Values.Server.File == configs.ServerFileQiniu {
		return doUploadBySimpleWithQiniu(ctx, dirs, file)
	}

	return doUploadBySimpleWithSystem(ctx, dirs, file)
}

func doUploadBySimpleWithSystem(ctx *gin.Context, dirs string, file *multipart.FileHeader) (error, *UploadBySimple) {

	filepath := "/upload"

	if !strings.HasPrefix(dirs, "/") {
		filepath += "/"
	}

	filepath += dirs

	if err := os.MkdirAll(config.Application.Runtime+filepath, 0750); err != nil {
		return err, nil
	}

	node, err := snowflake.NewNode(config.Values.Server.Node)
	if err != nil {
		return err, nil
	}

	generate := node.Generate()

	filename := generate.String() + path.Ext(file.Filename)

	filepath += "/" + filename

	err = ctx.SaveUploadedFile(file, config.Application.Runtime+filepath)
	if err != nil {
		return err, nil
	}

	return nil, &UploadBySimple{
		Name: filename,
		Path: filepath,
		Url:  config.Values.Server.Url + filepath,
	}
}

func doUploadBySimpleWithQiniu(ctx *gin.Context, dirs string, file *multipart.FileHeader) (error, *UploadBySimple) {

	policy := storage.PutPolicy{
		Scope: config.Values.Qiniu.Bucket,
	}

	mac := qbox.NewMac(config.Values.Qiniu.Access, config.Values.Qiniu.Secret)

	token := policy.UploadToken(mac)

	dump.P(token)

	resume := storage.NewFormUploader(nil)

	extra := storage.PutExtra{
		Params: map[string]string{
			"fileName": "",
		},
	}

	f, err := file.Open()
	if err != nil {
		return err, nil
	}

	var ret any

	err = resume.PutWithoutKey(ctx, &ret, token, bufio.NewReader(f), file.Size, &extra)
	if err != nil {
		dump.P(err)
		return err, nil
	}

	dump.P(ret)

	//node, err := snowflake.NewNode(config.Values.Server.Node)
	//if err != nil {
	//	return err, nil
	//}
	//
	//generate := node.Generate()
	//
	//filename := generate.String() + path.Ext(file.Filename)
	//
	//filepath += "/" + filename
	//
	//err = ctx.SaveUploadedFile(file, config.Application.Runtime+filepath)
	//if err != nil {
	return nil, nil
	//}
	//
	//return nil, &UploadBySimple{
	//	Name: filename,
	//	Path: filepath,
	//	Url:  config.Values.Server.Url + filepath,
	//}
}

type UploadBySimple struct {
	Name string
	Path string
	Url  string
}
