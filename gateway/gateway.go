package main

import (
	"context"
	"fmt"
	pb "web/protos/example.com/auth"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("starting")
	r := gin.Default()
	r.POST("/auth/req_pq", func(c *gin.Context) {
		conn, err := grpc.Dial("127.0.0.1:3313", grpc.WithInsecure())
		fmt.Println(conn)
		if err == nil {
			requestBody := pb.RequestPq{}
			c.BindJSON(&requestBody)
			//fmt.Println(requestBody)
			//nonce := c.Request.Header["Nonce"][0]
			//message_id, _ := strconv.Atoi(c.Request.Header["Message_id"][0])

			client := pb.NewWebClient(conn)
			answer, err := client.ReqPq(context.Background(), &requestBody)
			//answer, err := client.ReqPq(context.Background(), &pb.RequestPq{Nonce: nonce, MessageId: int32(message_id)})
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

		if conn != nil {
			defer conn.Close()
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
