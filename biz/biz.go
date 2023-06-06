package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "web/protos/example.com/biz"
)

type server struct {
	pb.UnimplementedBizServer
}

var db *sql.DB

func (s *server) GetUsers(ctx context.Context, in *pb.Request) (*pb.Result, error) {
	var query string
	if in.GetUserId() == 0 {
		log.Printf("Received: %v", in.GetMessageId())
		query = fmt.Sprintf("SELECT * FROM users limit 100", in.UserId)
	} else {
		log.Printf("Received: %v", in.GetMessageId())
		query = fmt.Sprintf("SELECT * FROM users where id = %d ", in.UserId)
	}
	return getUser(query, in.MessageId)
}
func (s *server) GetUsersWithSqlInject(ctx context.Context, in *pb.RequestSqlInject) (*pb.Result, error) {
	var query string
	if in.GetUserId() == "" {
		log.Printf("Received: %v", in.GetMessageId())
		query = fmt.Sprintf("SELECT * FROM users limit 100;")
	} else {
		log.Printf("Received: %v", in.GetMessageId())
		query = fmt.Sprintf("SELECT * FROM users where id = %s;", in.UserId)
	}
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
		err := users.Scan(&user.Name, &user.Family, &user.Age, &user.Sex, &user.Id, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		usersList = append(usersList, user)
	}
	return &pb.Result{Users: usersList, MessageId: messageId}, nil
}

func main() {
	var err error
	db, err = sql.Open("postgres", "postgres://baeldung:baeldung@localhost:5431/web?sslmode=disable")
	if err != nil {
		panic(err)
	}
	port := 3314
	lis, err2 := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err2 != nil {
		log.Fatalf("failed to listen: %v", err2)
	}
	s := grpc.NewServer()
	pb.RegisterBizServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
