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
		//"exp": time.Now().Add(time.Hour * 24).Unix(),
		//"exp": time.Now().Add(time.Hour * 1).Unix(),
		"exp": time.Now().Add(time.Minute * 30).Unix(),
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

	fmt.Println("tokenString:", tokenString)
	JwtString = tokenString

	fmt.Println("----------")
}

// JWTの検証
func JwtVerify(jwtToken string) (userid string, err error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			fmt.Println("アルゴリズムが正しくありません")
			return nil, nil
		}

		//return []byte("SECRET_KEY"), nil
		return []byte("prsk_key"), nil
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		fmt.Println("認証OK")
	} else {
		fmt.Println("認証エラー")
		return "", err
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		fmt.Println("not found user_id in " + jwtToken)
		return "", err
	}

	fmt.Println(userID)
	fmt.Println("----------")
	return userID, nil
}
