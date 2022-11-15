package tools

import (
	"bytes"
	"cloud_disk/greet/define"
	"context"
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jordan-wright/email"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// GetMd5 把密码进行加密处理
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// GenerateToken 生成用户的Token信息，使用用户ID，Identity和name生成，并且带上token的有效时间
func GenerateToken(id int, identity, name string, second int64) (string, error) {
	userClaim := define.UserClaim{
		Id:       id,
		Identity: identity,
		Name:     name,
		RegisteredClaims: jwt.RegisteredClaims{
			//设置token过期时间
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Second * time.Duration(second))},
		},
	}

	//设置签名的加密算法，并且生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	//加上我们自己的关键字生成一个自己的JWTToken
	tokenString, err := token.SignedString(define.JWTKey)
	if err != nil {
		return "", err
	}
	return tokenString, err

}

// AnalyseToken 解析token并且把解析出来的用户声明返回出去
func AnalyseToken(tokenString string) (*define.UserClaim, error) {
	//定义一个用户声明，用来返回token中解析出来的用户声明信息
	userClaim := new(define.UserClaim)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(*jwt.Token) (interface{}, error) {
		//这个方法属于引用的包中的方法，用来得到我们的key，所以这里我们直接使用匿名方法
		return define.JWTKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims.Valid {
		return nil, errors.New("token解析失败")
	}
	return userClaim, nil
}

// SendCode 通过邮箱发送验证码
func SendCode(userEmail string, code string) error {
	e := email.NewEmail() //new一个email对象
	//配置email的基本信息
	e.From = "Get <getUserEmail@163.com>" //用自己的邮箱
	e.To = []string{userEmail}
	e.Subject = "验证码已发送，请查收！"
	e.HTML = []byte("您的验证码：<b>" + code + "</b>")
	//设置163邮箱的中转协议，与邮箱密码（自己的邮箱密码）
	return e.SendWithTLS("smtp.163.com:465",
		smtp.PlainAuth("", "getUserEmail@163.com", define.MailPassword, "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
}

// GenerateCode 生成用来发送给用户邮箱中的验证码
func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	res := ""
	for i := 0; i < 6; i++ {
		res += strconv.Itoa(rand.Intn(10))
	}
	return res
}

// GetUUID 获取UUID作为我们用户的id使用
func GetUUID() string {
	u := uuid.New().String()
	if u == "" {
		log.Printf("生成uuid失败")
		return ""
	}
	return u
}

// TenCosUpload 从前端获取文件上传
func TenCosUpload(req *http.Request) (string, error) {
	u, _ := url.Parse(define.TenCosURL)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥和ID
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: os.Getenv(define.TenSecretID), //我是放在环境变量中获取的（自己可以去自己的腾讯云服务获取，上面注释有提示）
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: os.Getenv(define.TenSecretKey),
		},
	})

	//获取从前端传过来的文件(file类型的数据，key为file)
	file, header, err := req.FormFile("file")
	if err != nil {
		log.Printf("获取前端文件失败")
		return "", err
	}

	//指定上传到到那个文件夹和叫什么文件名（上传到存储桶中）,文件名使用uuid加文件类型（ path.Ext()可以获取文件后缀）
	key := "cloud_disk/" + GetUUID() + path.Ext(header.Filename)

	//通过流传输
	_, err = client.Object.Put(
		context.Background(), key, file, nil,
	)
	if err != nil {
		panic(err)
	}
	//返回上传文件的url
	return define.TenCosURL + "/" + key, nil
}

// InitTenCosPart 初始化tencentCOS分片上传
func InitTenCosPart(ext string) (string, string, error) {
	//设置上传到的文件夹和文件名
	key := "cloud_disk/" + GetUUID() + ext
	//TenCOsClient我们定义到initialize中，统一获取
	v, _, err := define.TenCosClient.Object.InitiateMultipartUpload(context.Background(), key, nil)
	if err != nil {
		log.Printf("初始化tencentCOS分片失败：%v", err)
		return "", "", err
	}

	//初始化获取UploadID
	UploadID := v.UploadID
	//返回两个数据一个key，另一个UploadID
	return key, UploadID, nil
}

// TenCosPartUpload 分片上传
func TenCosPartUpload(r *http.Request) (string, error) {
	////获取前端传过来的key,uploadId和分片编号partNumber
	key := r.PostForm.Get("key")
	UploadID := r.PostForm.Get("upload_id")
	partNumber, err := strconv.Atoi(r.PostForm.Get("part_number"))
	if err != nil {
		log.Printf("数据格式转换出错：%v", err)
		return "", nil
	}

	//开始读取文件
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("读取分片文件失败：%v", err)
		return "", err
	}

	/*
			这里如果和Put传输方式一样直接把multipart.File类型放进去的话，浏览器会报（411 missing content length错）
		可以采取官网的方式在InitiateMultipartUpload()初始化中传入opt参数，在opt参数结构中设定我们的contentLength属性值
		我们这里就使用buffer来处理
	*/
	buf := bytes.NewBuffer(nil)
	//把文件写到buffer流中
	io.Copy(buf, file)

	// opt可选
	resp, err := define.TenCosClient.Object.UploadPart( //上传分片
		context.Background(), key, UploadID, partNumber, bytes.NewReader(buf.Bytes()), nil,
	)
	if err != nil {
		return "", nil
	}
	//获取上传分块文件的md5值（ETag）
	PartETag := resp.Header.Get("ETag")
	//由于获取的ETag带有双引号，所以我们需要把他去掉
	return strings.Trim(PartETag, "\""), nil
}

// TenCosPartUploadComplete 文件分片上传成功
func TenCosPartUploadComplete(key, UploadID string, objects []cos.Object) error {
	//完成分片上传
	opt := &cos.CompleteMultipartUploadOptions{}
	//把所有分片PartNumber和ETag，用来校验块的准确性
	opt.Parts = append(opt.Parts, objects...)
	_, _, err := define.TenCosClient.Object.CompleteMultipartUpload(
		context.Background(), key, UploadID, opt,
	)
	if err != nil {
		return err
	}
	return nil
}
