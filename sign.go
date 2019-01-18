package express100api

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/astaxie/beego/logs"
)

//Sign -
type Sign struct {
	Values *url.Values
}

func (s *Sign) Sign(clientSecret string, signmethod string) string {
	plaintext := s.MakeSignPlainText(clientSecret)
	return s.MakeSign(plaintext, clientSecret, signmethod)
}

/**
*MakeSign
*生成签名
@param string srcstr 拼接签名源文字符串
@param string secretkey secretkey
@param string signmethod 加密方法
@return string 加密串
**/
func (s *Sign) MakeSign(srcstr, secretkey, signmethod string) string {
	var retstr string
	switch signmethod {
	case "HmacSHA1":
		retstr = base64.StdEncoding.EncodeToString([]byte(HmacSha1(srcstr, secretkey)))
	case "HmacSHA256":
		retstr = base64.StdEncoding.EncodeToString([]byte(HMacSHA256(srcstr, secretkey)))
	case "Md5":
		retstr = Md5string(srcstr)
		retstr = hex.EncodeToString([]byte(retstr))
		retstr = strings.ToUpper(retstr)
		logs.Info("%v", retstr)
	}
	return retstr
}

/**
*MakeSignPlainText
*生成拼接签名源字符串
@param values xhttp.URLEncoder 请求参数
@return plaintext string
**/
func (s *Sign) MakeSignPlainText(clientSecret string) string {

	//拼接参数
	paramstr := s.Buildparamstr()
	plaintext := clientSecret + paramstr + clientSecret

	logs.Info("%v", plaintext)

	return plaintext
}

/**
*buildparamstr
*拼接参数
@param values xhttp.URLEncoder 请求参数
@return string	返回拼接参数
**/
func (s *Sign) Buildparamstr() string {
	var paramstr string

	//请求数据按A~Z排序
	keys := make([]string, 0)
	for k, v := range (map[string][]string)(*s.Values) {
		if "" == k {
			continue
		}
		//空字符串
		if 0 == len(v) {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if key == "sign" {
			continue
		}
		paramstr = paramstr + key + fmt.Sprintf("%v", s.Values.Get(key))
	}

	return paramstr
}

func Md5string(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return string(h.Sum(nil))
}

//HmacSha1 -
func HmacSha1(data, key string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	return string(mac.Sum(nil))
}

//HMacSHA256 -
func HMacSHA256(data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	return string(mac.Sum(nil))
}
