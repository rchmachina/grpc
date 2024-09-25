package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	
	hashing "github.com/rchmachina/grpc/cmd/utils/hashing"
	auth "github.com/rchmachina/grpc/cmd/config/auth"

	//sg  "github.com/rchmachina/grpc/cmd/utils/setGetCtx"
	userPb "github.com/rchmachina/grpc/dto/authpb"
	"gorm.io/gorm"
)

type AuthService struct {
	userPb.UnimplementedAuthServiceServer
	Db *gorm.DB
	//FiberCtx  *fiber.App
}

func (a *AuthService) Login(ctx context.Context, req *userPb.LoginRequest) (*userPb.LoginResponse, error) {
	log.Println("isi user", req)
	var responseLogin *userPb.LoginResponse
	var result string
	paramsJSON, err := json.Marshal(map[string]interface{}{
		"email":          req.GetEmail(),
	})


	err = a.Db.Raw("SELECT * FROM users.get_user_by_email($1::jsonb)", string(paramsJSON)).Scan(&result).Error
	if err != nil {
		return responseLogin, err
	}

	// Unmarshal the result back into the responseLogin struct
	err = json.Unmarshal([]byte(result), &responseLogin)
	if err != nil {
		return responseLogin, err
	}

	expiredTime := time.Now().Add(4 * time.Hour).Unix()
	responseLogin.Expired = expiredTime
	token, errGenerateToken := auth.GenerateToken(responseLogin)
	if errGenerateToken != nil {
		log.Println("Error generating token:", errGenerateToken)
		return responseLogin, nil
	}
	log.Println("Generated token:", token)
	responseLogin.Token = token

	return responseLogin, nil
}

// // CreateUser method for registering a new user
func (a *AuthService) CreateUser(ctx context.Context, req *userPb.CreateUserRequest) (*userPb.CreateUserResponse, error) {
	log.Println("isi user", req)
	var responseCreate *userPb.CreateUserResponse
	var result string
	

	hashedPassword, err := hashing.HashingPassword(req.GetHashedPassword())
	if err != nil {
		return responseCreate, err
	}
	req.HashedPassword = hashedPassword

	paramsJSON, err := json.Marshal(req)

	err = a.Db.Raw("SELECT * FROM users.create_users($1::jsonb)", string(paramsJSON)).Scan(&result).Error
	if err != nil {
		return responseCreate, err
	}

	// Unmarshal the result back into the responseLogin struct
	err = json.Unmarshal([]byte(result), &responseCreate)
	if err != nil {
		return responseCreate, err
	}

	return responseCreate, nil

}
func (a *AuthService) TestingMw(ctx context.Context, req *userPb.Empty) (*userPb.ReturnTesting, error) {
	var response *userPb.ReturnTesting
	response = &userPb.ReturnTesting{GetTest: "success"}
	return response, nil

}
