package main

import (
	"context"
	"fmt"
	"log"

	// "fmt"
	"testing"

	pb "gRPC_GO_crud/proto"

	// "github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestingDatabaseConnection() *gorm.DB{
	host := "localhost"
	port := "5432"
	dbName := "crud_go_test"
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

func TestCreateUser(t *testing.T) {
	db := TestingDatabaseConnection()
	srv := &server{DB: db}
	req := &pb.CreateUserRequest{
		User: &pb.User{
			Name: "John Doe",
			Age:  30,
		},
	}
	ctx := context.Background()

	res, err := srv.CreateUser(ctx, req)

	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	if res == nil {
		t.Error("CreateUser response is nil")
	}

	db.Exec("TRUNCATE TABLE USERS")
}


func TestGetUser(t *testing.T) {
	db := TestingDatabaseConnection()
	srv := &server{DB: db}
	req := &pb.CreateUserRequest{
		User: &pb.User{
			Name: "John Doe",
			Age:  30,
		},
	}
	ctx := context.Background()
	res, err := srv.CreateUser(ctx, req)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	reqGet := &pb.ReadUserRequest{
		Id: res.User.Id,
	}
	resGet, errGet := srv.GetUser(ctx, reqGet)

	if errGet != nil {
		t.Errorf("Error getting user: %v", err)
	}
	if resGet == nil {
		t.Error("GetUser response is nil")
	}
	if resGet.User.Id != req.User.Id {
		t.Errorf("Expected user ID %s, got %s", req.User.Id, res.User.Id)
	}

	db.Exec("TRUNCATE TABLE USERS")
}

func TestUpdateUser(t *testing.T) {
	db := TestingDatabaseConnection()
	srv := &server{DB: db}
	req := &pb.CreateUserRequest{
		User: &pb.User{
			Name: "John Doe",
			Age:  30,
		},
	}
	ctx := context.Background()
	res, err := srv.CreateUser(ctx, req)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	reqUpdate := &pb.UpdateUserRequest{
		User: &pb.User{
			Id:   res.User.Id,
			Name: "Updated Name",
		},
	}

	resUpdate, errUpdate := srv.UpdateUser(ctx, reqUpdate)

	if errUpdate != nil {
		t.Errorf("Error updating user: %v", err)
	}
	if resUpdate.User.Name != reqUpdate.User.Name {
		t.Errorf("User details not updated properly")
	}

	db.Exec("TRUNCATE TABLE USERS")
}
