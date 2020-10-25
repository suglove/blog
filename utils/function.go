package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	redisClient "goblog/service/redis"
	"os"
	"strconv"
	"strings"
	"time"
)

func IndexForOne(i int, p, limit int64) int64 {
	s := strconv.Itoa(i)
	index, _ := strconv.ParseInt(s, 10, 64)
	return (p-1)*limit + index + 1
}

func IndexAddOne(i interface{}) int64 {
	index, _ := ToInt64(i)
	return index + 1
}

func IndexDecrOne(i interface{}) int64 {
	index, _ := ToInt64(i)
	return index - 1
}

func StringReplace(str, old, new string) string {
	return strings.Replace(str, old, new, -1)
}

func StringToTime(date interface{}) time.Time {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	ret, _ := time.ParseInLocation(timeLayout, date.(string), loc)
	return ret
}

func TimeStampToTime(timeStamp int32) time.Time {
	return time.Unix(int64(timeStamp), 0)
}

// ToInt64 类型转换，获得int64
func ToInt64(v interface{}) (re int64, err error) {
	switch v.(type) {
	case string:
		re, err = strconv.ParseInt(v.(string), 10, 64)
	case float64:
		re = int64(v.(float64))
	case float32:
		re = int64(v.(float32))
	case int64:
		re = v.(int64)
	case int32:
		re = v.(int64)
	default:
		err = errors.New("不能转换")
	}
	return
}

//图片上传
func UploadImg(Filename, img string) (folderPath, imgName string, err error) {

	//保存文件
	fileArr := strings.Split(Filename, ".")
	index := len(fileArr)
	str := fileArr[index-1]
	fileName := FileName()
	imgName = fileName + "." + str
	basePath := "static/uploads/"

	folderName := time.Now().Format("2006-01-02")
	folderPath = basePath + folderName
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步
		// 先创建文件夹
		os.Mkdir(folderPath, 0777)
		// 再修改权限
		os.Chmod(folderPath, 0777)
	}

	return
}

//随机生成MD5值文件名称
func FileName() (name string) {
	m := md5.New()
	m.Write([]byte(time.Now().String()))
	b := m.Sum(nil)
	name = hex.EncodeToString(b)
	return
}

//密码加密
func PasswdMd5(passwd string) (passwdMd5 string) {
	m := md5.New()
	salt := beego.AppConfig.String(("appname")) + passwd
	m.Write([]byte(salt))
	b := m.Sum(nil)
	passwdMd5 = hex.EncodeToString(b)
	return
}

//获取点击数
func GetClick(id int) int {
	redisKey := "articleClick"
	conn, err := redisClient.PoolConnect()
	defer conn.Close()
	if err == nil {
		id, _ = redis.Int(conn.Do("ZSCORE", redisKey, id))
	}
	return id
}
