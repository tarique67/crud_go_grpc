package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	pb "gRPC_GO_crud/proto"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	DatabaseConnection()
}


type User struct {
	ID        string `gorm:"primarykey"`
	Name      string
	Age       int32
	CreatedAt time.Time `gorm:"autoCreateTime:false"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
}

func DatabaseConnection() *gorm.DB{
	host := "localhost"
	port := "5432"
	dbName := "crud_go"
	dbUser := "postgres"
	password := "pass1234"
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB.AutoMigrate(User{})
	if err != nil {
		log.Fatal("Error connecting to the database...", err)
	}
	fmt.Println("Database connection successful...")
	return DB
}

var (
	port = flag.Int("port", 50051, "gRPC server port")
 )
  
 type server struct {
	DB *gorm.DB
	pb.UnimplementedUserServiceServer
 }

 func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	fmt.Println("Create User")
	user := req.GetUser()
	user.Id = uuid.New().String()
  
	data := User{
		ID:    user.GetId(),
		Name: user.GetName(),
		Age: user.GetAge(),
	}
  
	res := s.DB.Create(&data)
	if res.RowsAffected == 0 {
		return nil, errors.New("user creation unsuccessful")
	}
	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:    user.GetId(),
			Name: user.GetName(),
			Age: user.GetAge(),
		},
	}, nil
 }
  
 func (s *server) GetUser(ctx context.Context, req *pb.ReadUserRequest) (*pb.ReadUserResponse, error) {
	fmt.Println("Read User", req.GetId())
	var user User
	res := s.DB.Find(&user, "id = ?", req.GetId())
	if res.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return &pb.ReadUserResponse{
		User: &pb.User{
			Id:    user.ID,
			Name: user.Name,
			Age: user.Age,
		},
	}, nil
 }
  
 func (s *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	fmt.Println("Update User")
	var user User
	reqUser := req.GetUser()
  
	res := s.DB.Model(&user).Where("id=?", reqUser.Id).Updates(User{Name: reqUser.Name})
  
	if res.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
  
	return &pb.UpdateUserResponse{
		User: &pb.User{
			Id:    user.ID,
			Name: user.Name,
		},
	}, nil
 }

 func main() {
	fmt.Println("gRPC server running ...")
  
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
  
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
  
	s := grpc.NewServer()
    db := DatabaseConnection();
	pb.RegisterUserServiceServer(s, &server{DB: db})
  
	log.Printf("Server listening at %v", lis.Addr())
  
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
 }