package main

import (
	"crypto/tls"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const TokenExpireDuration = time.Hour * 24

var mySecret = []byte("秘密")
var cli = &http.Client{}
var failCount float32

type myClaims struct {
	UserID               string `json:"user_id"`
	jwt.RegisteredClaims        // 内嵌申明字段
}

const TotalRequest float32 = 190000

func main() {
	initClient()
	failCount = 0
	wg := &sync.WaitGroup{}
	wg.Add(4)
	beginTime := time.Now()
	go doBuy(0, 60000, wg)
	go doBuy(60001, 120000, wg)
	go doBuy(120001, 180000, wg)
	go doBuy(180001, 190000, wg)
	wg.Wait()
	endTime := time.Now()
	fmt.Printf("请求失败率：%f%%\n耗时: %v \n失败次数: %d", failCount*100/TotalRequest, endTime.Sub(beginTime), int(failCount))
}

func initClient() {
	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        5000,              // 空闲(保持活跃)连接的最大数量
		MaxConnsPerHost:     5000,              // 限制每台主机的连接总数
		MaxIdleConnsPerHost: 5000,              // 保持的最大空闲连接数
		IdleConnTimeout:     300 * time.Second, // 空闲连接(保持连接)在关闭自己之前保持空闲的最大时间
	}
	cli = &http.Client{
		Transport: tr,
		Timeout:   9 * time.Second,
	}
}

func doBuy(begin int, end int, group *sync.WaitGroup) {
	defer group.Done()
	// 同步组
	g := &sync.WaitGroup{}
	g.Add(end - begin + 1)
	// 并发请求
	for i := begin; i <= end; i++ {
		token, err := GenToken(strconv.Itoa(i))
		if err != nil {
			panic("生成token失败：" + err.Error())
		}
		time.Sleep(time.Nanosecond)
		go func(*sync.WaitGroup, string) {
			if err := sendRequest(g, token); err != nil {
				l := &sync.Mutex{}
				l.Lock()
				failCount++
				l.Unlock()
				fmt.Println("请求出错: ", err.Error())
			}
		}(g, token)
	}
	g.Wait()
}

func sendRequest(wg *sync.WaitGroup, token string) error {
	// 1. 设置请求参数
	params := url.Values{}
	params.Set("token", token)
	params.Set("ticket_id", "h1k4J7Dyt0")
	params.Set("count", "1")
	// 2. 伪造请求
	req, err := http.NewRequest("POST", "http://localhost:8080/buy", strings.NewReader(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))
	// 3. 发送请求
	resp, err := cli.Do(req)
	defer wg.Done()
	if err != nil {
		return err
	}
	_, _ = ioutil.ReadAll(resp.Body)
	return nil
}

// GenToken 生成JWT
func GenToken(userID string) (string, error) {
	// 创建一个我们自己的声明
	claims := myClaims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "my-project", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}
