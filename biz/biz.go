package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	pb "web/protos/example.com/biz"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedWebServer
}

var db *sql.DB

func (s *server) RequestBiz(ctx context.Context, in *pb.Request) (*pb.Result, error) {
	log.Printf("Received: %v", in.GetMessageId())
	query := fmt.Sprintf("SELECT * FROM users where id = %d ", in.UserId)
	return getUser(query, in.MessageId)
}
func (s *server) RequestBizSqlInject(ctx context.Context, in *pb.RequestSqlInject) (*pb.Result, error) {
	log.Printf("Received: %v", in.GetMessageId())
	query := fmt.Sprintf("SELECT * FROM users where id = %s ", in.UserId)
	return getUser(query, in.MessageId)
}
func getUser(query string, messageId int32) (*pb.Result, error) {
	user := new(pb.User)
	var usersList []*pb.User
	users, err := db.Query(query)
	if err != nil {
		log.Fatalf("can't select from db")
	}
	for users.Next() {
		err := users.Scan(&user)
		if err != nil {
			return nil, err
		}
		usersList = append(usersList, user)
	}
	return &pb.Result{Users: usersList, MessageId: messageId}, nil
}

func main() {
	db, _ = sql.Open("postgres", "postgres://baeldung:baeldung@localhost:5431/web")

	port := 3314
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
