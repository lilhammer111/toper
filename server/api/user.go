package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
	"to-persist/server/form"
	"to-persist/server/global"
	"to-persist/server/model"
)

func GetUserList(c *gin.Context) {

}

func Login(c *gin.Context) {
	loginForm := form.LoginForm{}
	err := c.ShouldBind(&loginForm)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	user := model.User{}
	// check password
	if res := global.MysqlDB.Where("username = ?", loginForm.Name).First(&user); res.RowsAffected == 0 {
		c.Status(http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginForm.Password))
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken(user.Model.ID)
	if err != nil {
		zap.S().Error("failed to generate token, because ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}

func Register(c *gin.Context) {
	//
	registerForm := form.RegisterForm{}
	err := c.ShouldBind(&registerForm)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",
			global.Config.RedisConfig.AddrConfig.Host,
			global.Config.RedisConfig.AddrConfig.Port,
		),
	})

	rightSmsCode, err := rdb.Get(context.Background(), registerForm.Mobile).Result()

	if err == redis.Nil || registerForm.SMSCode != rightSmsCode {
		c.Status(http.StatusBadRequest)
		return
	}

	user := model.User{}

	if res := global.MysqlDB.Where("mobile = ?", registerForm.Mobile).First(&user); res.RowsAffected != 0 {
		c.Status(http.StatusConflict)
		return
	}

	//Username: registerForm.Name,
	//	Mobile:   registerForm.Mobile,
	//		Password: registerForm.Password,
	user.Username = registerForm.Name
	user.Mobile = registerForm.Mobile

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerForm.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	if res := global.MysqlDB.Create(&user); res.RowsAffected == 0 {
		zap.S().Error("res err:", res.Error.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	token, err := GenerateToken(user.Model.ID)
	if err != nil {
		zap.S().Error("failed to generate token, because ", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
	})

}

func GenerateToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Issuer:    "Toper",
		Subject:   strconv.Itoa(int(userID)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(global.Config.JwtConfig.JwtKey))
}
