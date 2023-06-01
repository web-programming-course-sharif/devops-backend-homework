package web

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	pb "web/example.com/biz"
)

var db *sql.DB

func (s *server) RequestBiz(ctx context.Context, in *pb.Request) (*pb.Result, error) {
	log.Printf("Received: %v", in.GetMessageId())
	query := fmt.Sprintf("SELECT * FROM users where id = %d ", in.UserId)
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
	return &pb.Result{Users: usersList, MessageId: in.GetMessageId()}, nil
}
func main() {
	db, _ = sql.Open("postgres", "postgres://user:pass@localhost/db")

}
