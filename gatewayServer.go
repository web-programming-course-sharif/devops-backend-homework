package web

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.POST("/auth/req_pq", func(c *gin.Context) {
		var opts []grpc.DialOption
		conn, err := grpc.Dial("localhost:313", opts...)
		if err != nil {
			nonce := c.Request.Header["nonce"]
			message_id := c.Request.Header["message_id"]

			client := pb.NewRouteGuideClient(conn)
			answer, err := client.RequestPq(context.Background(), pb.RequestPq{nonce: nonce, message_id: message_id})
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

	r.GET("/auth/req_DH_params", func(c *gin.Context) {
		var opts []grpc.DialOption
		conn, err := grpc.Dial("localhost:313", opts...)
		if err != nil {
			nonce := c.Request.Header["nonce"]
			server_nonce := c.Request.Header["server_nonce"]
			message_id := c.Request.Header["message_id"]
			a := c.Request.Header["a"]

			client := pb.NewRouteGuideClient(conn)
			answer, err := client.RequestDh(context.Background(), pb.RequestPq{nonce: nonce, server_nonce: server_nonce,
				message_id: message_id, a: a})
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
