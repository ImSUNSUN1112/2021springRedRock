package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

/*
一个用于签发token,一个用于验证token与可能的更新
 */

type Header struct {
	//加密算法
	Alg string `json:"alg"`
	//类别
	Typ string `json:"typ"`
}

type Payload struct {
	Sub Sub `json:"sub"`
	Iss string `json:"iss"`
	Exp string `json:"exp"`
	Iat string `json:"iat"`
}

type Sub struct {
	Username string `json:"username"`
	Telephone string `json:"telephone"`
}

type JWT struct {
	Header    Header
	Payload   Payload
	Signature string
	Token     string
}

func GetJWT(username string,telephone string) JWT {

	var jwt JWT

	fmt.Println("将要进行封装的sub是"+username+" "+telephone+" ")

	jwt.Header = Header{
		Typ: "JWT",
		Alg: "HS256",
	}
	jwt.Payload = Payload{
		Iss: "sunsun",
		Exp: strconv.FormatInt(time.Now().Add(3*time.Hour).Unix(), 10),
		Iat: strconv.FormatInt(time.Now().Unix(), 10),
		Sub: Sub{
			Username: username,
			Telephone: telephone,
		},
	}

	h, _ := json.Marshal(jwt.Header)

	p, _ := json.Marshal(jwt.Payload)

	baseH := base64.StdEncoding.EncodeToString(h)

	baseP := base64.StdEncoding.EncodeToString(p)

	secret := baseH + "." + baseP

	key := "sunsun"

	//sha256加密
	mac := hmac.New(sha256.New, []byte(key))

	mac.Write([]byte(secret))

	s := mac.Sum(nil)

	jwt.Signature = base64.StdEncoding.EncodeToString(s)

	jwt.Token = secret + "." + jwt.Signature

	fmt.Println("token是"+jwt.Token)

	return jwt
}

//检查token
func CheckJWT(token string)(jwt JWT,err error){
	err = errors.New("token error")

	//划分token
	arr := strings.Split(token,".")

	//检测token字段长度
	if len(arr)<3{
		fmt.Println("err="+err.Error())
		return jwt,err
	}

	//head部分转化
	baseH:=arr[0]
	h,err:=base64.StdEncoding.DecodeString(baseH)
	if err != nil {
		fmt.Println("decode header", err)
		return
	}
	err = json.Unmarshal(h,&jwt.Header)
	if err != nil {
		fmt.Println("unmarshal header", err)
		return
	}

	//playload部分转化
	baseP := arr[1]
	p, err := base64.StdEncoding.DecodeString(baseP)
	if err != nil {
		fmt.Println("decode payload", err)
		return
	}
	err = json.Unmarshal(p, &jwt.Payload)
	if err != nil {
		fmt.Println("unmarshal payload", err)
		return
	}

	//检测过期时间
	exp,err:=strconv.Atoi(jwt.Payload.Exp)
	if err!=nil {
		return
	}
	//签发时间已过
	if int64(exp)>time.Now().Unix(){
		return jwt,errors.New("time out")
	}


	//secret部分转化
	bases := arr[2]
	s1, err := base64.StdEncoding.DecodeString(bases)
	if err != nil {
		fmt.Println("decode secret", err)
		return
	}

	//公开部分再加密，与加密部分比较
	se := baseH + "." + baseP
	w:=[]byte("sunsun")
	mac:=hmac.New(sha256.New,w)
	mac.Write([]byte(se))
	s2:=mac.Sum(nil)

	if string(s1)!=string(s2){
		return jwt,err
	}else {
		//验证通过,返回一个新签发的jwt与nil
		return GetJWT(jwt.Payload.Sub.Username,jwt.Payload.Sub.Telephone),nil
	}

}
