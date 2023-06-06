package main

import (
	"context"
	"crypto/sha1"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"time"
	pb "web/protos/example.com/auth"
	"web/redis"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedWebServer
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

func (s *server) ReqPq(ctx context.Context, in *pb.RequestPq) (*pb.ResultPq, error) {
	log.Printf("Received: %v", in.GetMessageId())
	serverNonce := randomString(20)
	key := []byte(serverNonce + in.GetNonce())
	sha1Hash := fmt.Sprintf("%x", sha1.Sum(key))
	g := rand.Intn(8) + 2
	redis.Rdb.Set(ctx, sha1Hash, g, time.Minute*20)
	println(serverNonce)
	return &pb.ResultPq{Nonce: in.GetNonce(), ServerNonce: serverNonce, MessageId: in.GetMessageId() + 1, P: 23, G: int32(g)}, nil
}
func (s *server) Req_DHParams(ctx context.Context, in *pb.Request_DH) (*pb.Result_DH, error) {
	log.Printf("Received: %v", in.GetMessageId())
	key := []byte(in.GetServerNonce() + in.GetNonce())
	sha1Hash := fmt.Sprintf("%x", sha1.Sum(key))
	g, _ := redis.Rdb.Get(ctx, sha1Hash).Int()
	b := rand.Intn(8) + 2
	l := (in.GetA() ^ int32(b)) % 23
	redis.Rdb.Set(ctx, sha1Hash, l, time.Hour*24*365)
	return &pb.Result_DH{Nonce: in.Nonce, ServerNonce: in.GetServerNonce(), MessageId: in.GetMessageId() + 1, B: int32(g) ^ (l)}, nil
}

func main() {
	redis.ConnectToRedis()
	port := 3313
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterWebServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
