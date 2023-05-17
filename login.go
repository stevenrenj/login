package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"login/controller"
	"login/models"
	"math/rand"
	"net/http"
)

func RandmonString(n int) string {
	var letters = []byte("asdfghjklqwertyuiopzxcvbnmASDFGHJKLQWERTYUIOPZXCVBNM")
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user models.User
	db.Debug().Where("telephone=?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

func main() {
	db := controller.InitDB()
	defer db.Close()
	r := gin.Default()
	r.POST("/api/auth/register", func(c *gin.Context) {
		// 获取参数
		name := c.PostForm("name")
		password := c.PostForm("password")
		telephone := c.PostForm("telephone")
		// 数据验证
		if len(telephone) != 11 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
			return
		}
		if len(password) < 6 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
			return

		}
		// 如果名称没有传，给一个10位的随机字符串
		if len(name) == 0 {
			name = RandmonString(10)
		}
		log.Println(name, password, telephone)
		// 判断手机号是否存在
		if isTelephoneExist(db, telephone) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号已存在"})
			return
		}
		// 创建用户
		newUser := models.User{
			Name:      name,
			Password:  password,
			Telephone: telephone,
		}
		db.Create(&newUser)
		// 返回结果

		c.JSON(200, gin.H{
			"msg": "注册成功",
		})
	})
	panic(r.Run()) // 监听并在 0.0.0.0:8080 上启动服务
}
