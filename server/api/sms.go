package api

import (
	"context"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"to-persist/server/global"
)

func GenerateSmsCode(width int) string {
	//生成width长度的短信验证码

	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)

	var sb strings.Builder
	for i := 0; i < width; i++ {
		_, _ = fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func SendSms(c *gin.Context) {
	mobile := c.Query("mobile")

	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: &global.AccessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: &global.AccessKeySecret,
		// 访问的域名
	}
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")

	client, _ := dysmsapi.NewClient(config)

	smsCode := GenerateSmsCode(6)
	//smsCode := "1234"

	request := &dysmsapi.SendSmsRequest{}
	request.SetPhoneNumbers(mobile)
	request.SetSignName("阿里云短信测试")
	request.SetTemplateCode("SMS_154950909")
	request.SetTemplateParam(fmt.Sprintf("{\"code\":%s}", smsCode))

	response, err := client.SendSms(request)
	if err != nil {
		zap.S().Error("failed to send sms, because ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if *response.StatusCode != http.StatusOK {
		zap.S().Errorf("failed to send sms, because %s", *response.Body.Message)
		c.Status(http.StatusInternalServerError)
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			global.Config.RedisConfig.AddrConfig.Host,
			global.Config.RedisConfig.AddrConfig.Port,
		),
	})
	rdb.Set(context.Background(), mobile, smsCode, time.Duration(global.Config.RedisConfig.Expire)*time.Second)

	c.Status(http.StatusOK)
}
