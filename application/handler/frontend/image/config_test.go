package image

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/admpub/godotenv"
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/s3manager/s3client"
	"github.com/admpub/pp"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
)

func dump(a ...interface{}) {
	pp.Println(a...)
}

func TestStat(t *testing.T) {
	projectDir := filepath.Join(echo.Wd(), `../../../../`)
	envFile := filepath.Join(projectDir, `.env`)
	err := godotenv.Overload(envFile)
	if err != nil {
		panic(err)
	}
	cfg := &dbschema.NgingCloudStorage{
		Key:      os.Getenv(`S3_KEY`),
		Secret:   os.Getenv(`S3_SECRET`),
		Secure:   `Y`,
		Region:   os.Getenv(`S3_REGION`),
		Bucket:   os.Getenv(`S3_BUCKET`),
		Endpoint: os.Getenv(`S3_ENDPOINT`),
		Baseurl:  os.Getenv(`S3_BASEURL`),
	}
	//echo.Dump(cfg)
	mgr := s3client.New(cfg, 1024000)
	_, err = mgr.Connect()
	if err != nil {
		panic(err)
	}

	// =========================
	// 以Exists方式判断文件是否存在
	// =========================

	// 测试不存在的文件
	exists, err := mgr.Exists(context.Background(), `/nofile.f`)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, false, exists)

	// 测试存在的文件
	filePath := `/public/upload/film/2021/11/14/169684569803456512.jpg`
	exists, err = mgr.Exists(context.Background(), filePath)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, true, exists)

	// =========================
	// 以 Get 方式判断文件是否存在
	// =========================

	// 测试存在的文件
	obj, err := mgr.Get(context.Background(), filePath)
	if err != nil {
		panic(err)
	}
	dump(obj)
	exists, err = mgr.StatIsExists(obj.Stat())
	if err != nil {
		panic(err)
	}
	assert.Equal(t, true, exists)

	// 测试不存在的文件
	obj, err = mgr.Get(context.Background(), `/nofile.f`)
	if err != nil {
		panic(err)
	}
	dump(obj)
	exists, err = mgr.StatIsExists(obj.Stat())
	if err != nil {
		panic(err)
	}
	assert.Equal(t, false, exists)
}
