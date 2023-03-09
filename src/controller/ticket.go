package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"robTickets/src/logic"
	"robTickets/src/util"
	"strconv"
)

//func AuthHandler(c *gin.Context) {
//	// 用户发过来用户名和密码
//	var user models.User
//	if err := c.ShouldBindJSON(&user); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": 2001,
//			"msg":  "参数无效",
//		})
//		return
//	}
//
//	// 检验用户名和密码是否正确
//	if user.UserID == "nxj" && user.Password == "123123" {
//		// 生成token
//		tokenString, err := util.GenToken(user.UserID)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{
//				"code": 2001,
//				"mas":  "服务繁忙",
//			})
//			return
//		}
//		c.JSON(http.StatusOK, gin.H{
//			"code": 2000,
//			"msg":  "success",
//			"data": gin.H{"token": tokenString},
//		})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"code": 2002,
//		"msg":  "鉴权失败",
//	})
//	return
//}

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.PostForm("token")
		if token == "" {
			c.JSON(http.StatusOK, "需要登录")
			c.Abort()
			return
		}

		mc, err := util.ParseToken(token)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效token",
			})
			c.Abort()
			return
		}
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set("user_id", mc.UserID)
		// 后续的处理函数可以用过c.Get("userID")来获取当前请求的用户信息
		c.Next()
	}
}

// BuyTicket 买票
func BuyTicket(c *gin.Context) {
	userID := c.GetString("user_id")
	ticketID := c.PostForm("ticket_id")
	number := c.PostForm("count")

	quantity, _ := strconv.Atoi(number)

	// 更新票数量
	order, ok, err := logic.BuyTicket(userID, ticketID, quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "服务器繁忙"})
		log.Println("抢票异常: ", err.Error())
		return
	}
	if ok {
		c.JSON(http.StatusOK, gin.H{"msg": "购买成功", "data": order})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "购买失败"})
	return
}

// PayOrder 支付订单
func PayOrder(c *gin.Context) {

}
