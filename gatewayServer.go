package main

import (
	"fmt"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	pb "web/example.com/web"
	"strconv"
)

func main() {
	fmt.Println("starting")
	r := gin.Default()
	r.POST("/auth/req_pq", func(c *gin.Context) {
		conn, err := grpc.Dial("127.0.0.1:3313")
		fmt.Println(conn)
		if err != nil && conn!=nil {
			//requestBody := pb.RequestPq{}
			//c.BindJSON(&requestBody)
			//fmt.Println(requestBody)
			nonce := c.Request.Header["Nonce"][0]
			message_id,_ := strconv.Atoi(c.Request.Header["Message_id"][0])

			client := pb.NewWebClient(conn)
			//answer, err := client.ReqPq(context.Background(), &requestBody)
			answer, err := client.ReqPq(context.Background(), &pb.RequestPq{Nonce:nonce, MessageId: int32(message_id)})
			if err != nil {
				c.JSON(200, gin.H{
					"answer": answer,
				})
			} else {
				c.JSON(403, gin.H{
					"answer": "Can not authenticate.",
				})
			}
		} else {
			c.JSON(403, gin.H{
				"answer": "Can not connect.",
			})
		}

		if conn != nil{
			defer conn.Close()
		}
	})

	r.POST("/auth/req_DH_params", func(c *gin.Context) {
		//var opts []grpc.DialOption
		conn, err := grpc.Dial("127.0.0.1:3313")
		if err != nil {
			nonce := c.Request.Header["nonce"][0]
			server_nonce := c.Request.Header["server_nonce"][0]
			message_id,_ := strconv.Atoi(c.Request.Header["message_id"][0])
			a,_ := strconv.Atoi(c.Request.Header["a"][0])

			client := pb.NewWebClient(conn)
			answer, err := client.Req_DHParams(context.Background(), &pb.Request_DH{Nonce: nonce, ServerNonce: server_nonce,
				MessageId: int32(message_id), A: int32(a)})
			if err != nil {
				c.JSON(200, gin.H{
					"answer": answer,
				})
			} else {
				c.JSON(403, gin.H{
					"answer": "Can not authenticate.",
				})
			}
		} else {
			c.JSON(403, gin.H{
				"answer": "Can not authenticate.",
			})
		}

		defer conn.Close()
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
