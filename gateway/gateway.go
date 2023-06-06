package main

import (
	"context"
	"fmt"
	"github.com/JGLTechnologies/gin-rate-limit"
	"sync"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net/http"
	"time"
	pba "web/protos/example.com/auth"
	pbb "web/protos/example.com/biz"
)

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	fmt.Println(info.ResetTime)
	fmt.Println(info)
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

var (
	requestCount         int
	requestCountMux      sync.Mutex
	userBlockTimes       map[string]time.Time
	userBlockTimesMux    sync.Mutex
	maxRequestsPerSecond = 99
	blockDuration        = time.Hour * 24
)

func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		requestCountMux.Lock()
		requestCount++
		count := requestCount
		requestCountMux.Unlock()

		userBlockTimesMux.Lock()
		if userBlockTimes == nil {
			userBlockTimes = make(map[string]time.Time)
		}
		blockTime, exists := userBlockTimes[ip]
		userBlockTimesMux.Unlock()

		if count > maxRequestsPerSecond {
			if exists && time.Since(blockTime) < blockDuration {
				c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
					"message": "You are blocked for 24 hours.",
				})
				return
			}

			userBlockTimesMux.Lock()
			userBlockTimes[ip] = time.Now()
			userBlockTimesMux.Unlock()
		}

		c.Next()
	}
}

func main() {
	fmt.Println("starting")
	r := gin.Default()
	rateLimit := rateLimitMiddleware()
	connPba, errPba := grpc.Dial("127.0.0.1:3313", grpc.WithInsecure())
	connPbb, errPbb := grpc.Dial("127.0.0.1:3314", grpc.WithInsecure())
	if errPbb != nil || errPba != nil || connPba == nil || connPbb == nil {
		panic("err")
	}
	r.POST("/auth/req_pq", rateLimit, func(c *gin.Context) {
		requestBody := pba.RequestPq{}
		c.BindJSON(&requestBody)

		client := pba.NewWebClient(connPba)
		answer, err := client.ReqPq(context.Background(), &requestBody)
		//answer, err := client.ReqPq(context.Background(), &pba.RequestPq{Nonce: nonce, MessageId: int32(message_id)})
		if err != nil {
			c.JSON(403, gin.H{
				"answer": "Can not authenticate.",
			})
		} else {
			c.JSON(200, gin.H{
				"answer": answer,
			})
		}
	})

	r.POST("/auth/req_DH_pq", rateLimit, func(c *gin.Context) {
		requestBody := pba.Request_DH{}
		c.BindJSON(&requestBody)

		client := pba.NewWebClient(connPba)
		answer, err := client.Req_DHParams(context.Background(), &requestBody)
		//answer, err := client.ReqPq(context.Background(), &pba.RequestPq{Nonce: nonce, MessageId: int32(message_id)})
		if err != nil {
			c.JSON(403, gin.H{
				"answer": "Can not authenticate.",
			})
		} else {
			c.JSON(200, gin.H{
				"answer": answer,
			})
		}
	})

	r.POST("/biz/get_users_with_sql_inject", func(c *gin.Context) {
		requestBody := pbb.RequestSqlInject{}
		c.BindJSON(&requestBody)

		client := pbb.NewBizClient(connPbb)
		answer, err := client.GetUsersWithSqlInject(context.Background(), &requestBody)
		//answer, err := client.ReqPq(context.Background(), &pba.RequestPq{Nonce: nonce, MessageId: int32(message_id)})
		if err != nil {
			c.JSON(403, gin.H{
				"answer": "Can not authenticate.",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"answer": answer,
			})
			c.Writer.WriteHeader(http.StatusOK)
		}
	})

	r.POST("/biz/get_users", func(c *gin.Context) {
		requestBody := pbb.Request{}
		c.BindJSON(&requestBody)

		client := pbb.NewBizClient(connPbb)
		answer, err := client.GetUsers(context.Background(), &requestBody)
		//answer, err := client.ReqPq(context.Background(), &pba.RequestPq{Nonce: nonce, MessageId: int32(message_id)})
		if err != nil {
			c.JSON(403, gin.H{
				"answer": "Can not authenticate.",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"answer": answer,
			})
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
