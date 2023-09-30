package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var JwtString string

// JWTの生成
func JwtGenerater(userid string) {
	// Claimsオブジェクトの作成
	claims := jwt.MapClaims{
		"user_id": userid,
		//"password": "12345",
		//"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}
	// ヘッダーとペイロードの生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Header: map[string]interface {}{"alg":"HS256", "typ":"JWT"}
	fmt.Printf("Header: %#v\n", token.Header)
	// Claims: jwt.MapClaims{"exp":1634051243, "user_id":12345678}
	fmt.Printf("Claims: %#v\n", token.Claims)

	// トークンに署名を付与
	//tokenString, _ := token.SignedString([]byte("SECRET_KEY"))
	tokenString, _ := token.SignedString([]byte("prsk_key"))
	//
	fmt.Println("tokenString:", tokenString)
	JwtString = tokenString
	//
	fmt.Println("----------")
}

// JWTの検証
func JwtVerify(jwtToken string) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		//return []byte("SECRET_KEY"), nil
		return []byte("prsk_key"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//fmt.Printf("user_id: %v\n", int64(claims["user_id"].(float64)))
		fmt.Printf("user_id: %v\n", string(claims["user_id"].(string)))
		//fmt.Printf("password: %v\n", string(claims["password"].(string)))
		fmt.Printf("exp: %v\n", int64(claims["exp"].(float64)))
		fmt.Printf("alg: %v\n", token.Header["alg"])
		//
		//
		fmt.Println("----------")
		// unix時刻を日時に変換し有効時間を確認
		//unix := int64(claims["exp"].(int64))
		unix := int64(claims["exp"].(float64))
		dtFromUnix := time.Unix(unix, 0)
		//
		//nowtime := time.Now()
		nowtime := time.Now().Add(time.Hour * 2)
		// 時刻比較、現在時刻を超えているか、超えていたらfalse
		diff := dtFromUnix.After(nowtime)
		//
		fmt.Println(nowtime)
		fmt.Println(dtFromUnix)
		fmt.Println(diff)
		//
		fmt.Println("----------")
		nowunix := time.Now().Unix()
		//nowunix := time.Now().Add(time.Hour * 2).Unix()
		fmt.Println("制限時刻:", unix)
		fmt.Println("現在時刻:", nowunix)
		if unix < nowunix {
			fmt.Println("制限時刻を超えています")
		} else {
			fmt.Println("制限時刻を超えていません")
		}
		//
		//
	} else {
		fmt.Println("認証エラー")
		fmt.Println(err)
	}
	fmt.Println("----------")
}
