package test

import (
	"fmt"
	"os"
	"testing"
)

//获取repository_pool中path的key值
func TestSubString(t *testing.T) {
	str := "https://examplebucket-1250000000.cos.ap-guangzhou.myqcloud.com/cloud_disk/test.mp4"
	str1 := "https://examplebucket-1250000000.cos.ap-guangzhou.myqcloud.com/"
	fmt.Println(len(str))
	fmt.Println(len(str1))
	str2 := str[len(str1):]
	fmt.Println(str2)
	//输出
	//82
	//63
	//cloud_disk/test.mp4

	//获取当前文件路径
	fmt.Println(os.Getwd())
}
