package test

import (
	"bytes"
	"cloud_disk/greet/define"
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
	"testing"
)

//测试Upload方式上传文件
func TestFileUpload(t *testing.T) {
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶region可以在COS控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse("https://examplebucket-1250000000.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: os.Getenv(define.TenSecretID), //我是放在魂晶变量中获取的（自己可以去自己的腾讯云服务获取，上面注释有提示）
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: os.Getenv(define.TenSecretKey),
		},
	})

	//指定上传到到那个文件夹和叫什么文件名（上传到存储桶中）
	key := "cloud_disk/test.jpg"

	_, _, err := client.Object.Upload(
		context.Background(), key, "./img/test.jpg", nil,
	)
	if err != nil {
		t.Fatal(err)
	}
}

//测试Put方式上传文件
func TestFileUploadByPut(t *testing.T) {
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶region可以在COS控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse("https://examplebucket-1250000000.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: os.Getenv(define.TenSecretID), //我是放在魂晶变量中获取的（自己可以去自己的腾讯云服务获取，上面注释有提示）
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: os.Getenv(define.TenSecretKey),
		},
	})

	// Case1 使用 Put 上传对象
	key := "cloud_disk/test.jpg"
	f, err := os.ReadFile("./img/test.jpg")
	//opt := &cos.ObjectPutOptions{
	//	ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
	//		ContentType: "text/html",
	//	},
	//	ACLHeaderOptions: &cos.ACLHeaderOptions{
	//		// 如果不是必要操作，建议上传文件时不要给单个文件设置权限，避免达到限制。若不设置默认继承桶的权限。
	//		XCosACL: "private",
	//	},
	//}
	_, err = client.Object.Put(context.Background(), key, bytes.NewReader(f), nil)
	if err != nil {
		panic(err)
	}

	// Case 2 使用 PUtFromFile 上传本地文件到COS
	//filepath := "./test"
	//_, err = client.Object.PutFromFile(context.Background(), key, filepath, opt)
	//if err != nil {
	//	panic(err)
	//}

	// Case 3 上传 0 字节文件, 设置输入流长度为 0
	//_, err = client.Object.Put(context.Background(), key, strings.NewReader(""), nil)
	//if err != nil {
	//	// ERROR
	//}
}
