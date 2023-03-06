package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func TestRomanTokne(t *testing.T) {
	token := GengerateToken(99)
	fmt.Println(len(token))
	fmt.Println(token)
}

func TestModToken(t *testing.T) {
	token1 := Token{TokenString: GengerateToken(10), GenerateTime: GengerateToken(20)}
	ModTokenFile(token1, "token1")
	token2, _ := GetToken("token1")
	fmt.Println(token1, token2)
}

func TestS3(t *testing.T) {

	access_key := "xxxxxxxxxxxxx"
	secret_key := "xxxxxxxxxxxxx"
	end_point := "https://cos.ap-beijing.myqcloud.com" //endpoint设置，不要动

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(access_key, secret_key, ""),
		Endpoint:    aws.String(end_point),
		Region:      aws.String("ap-beijing"),
	})
	if err != nil {
		log.Println("Error creating session: ", err)
	}
	str1, err := get(sess, "恐咖兵糖的开始页.md") // retrun []byte
	fmt.Println(string(str1))
	fmt.Println(daily_file_key())
	fmt.Println(timeFmt("15:04"))
	list := get_3_daily_list(sess)
	for i, v := range list {
		fmt.Print(i)
		fmt.Println(v)
	}

	if err != nil {
		fmt.Println(err)
	}
}
