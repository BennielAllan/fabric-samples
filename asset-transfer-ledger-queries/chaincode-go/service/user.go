package service

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/hyperledger/fabric-samples/asset-transfer-ledger-queries/chaincode-go/utils"
	"time"
)

// User 用户表
type User struct {
	DocType      string  `json:"docType"`
	ID           string  `json:"id"` // 身份证ID
	Username     string  `json:"username"`
	PasswordHash string  `json:"password_hash"`
	Email        string  `json:"email"`
	Role         string  `json:"role"`
	Balance      float64 `json:"balance"`
	CreatedAt    int64   `json:"created_at"`
	UpdatedAt    int64   `json:"updated_at"`
}

// Register 注册
func (t *SimpleChaincode) Register(ctx contractapi.TransactionContextInterface,
	ID, UserName, password, email string) error {
	user := &User{
		DocType:      "user",
		ID:           ID,
		Username:     UserName,
		PasswordHash: utils.Hash(password),
		Email:        email,
		Role:         "user",
		Balance:      1000,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}
	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user info: %v", err)
	}
	if err = ctx.GetStub().PutState(user.ID, userBytes); err != nil {
		return fmt.Errorf("failed to PutState userBytes %s: %v", userBytes, err)
	}
	return nil
}

// Login 登录
func (t *SimpleChaincode) Login(ctx contractapi.TransactionContextInterface, userID string, password string) error {
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return fmt.Errorf("failed to get user info: %v", err)
	}
	if userBytes == nil {
		return fmt.Errorf("user info not found")
	}
	user := &User{}
	if err = json.Unmarshal(userBytes, &user); err != nil {
		return fmt.Errorf("failed to unmarshal user info: %v", err)
	}
	if user.PasswordHash == utils.Hash(password) {
		return nil
	}
	return nil
}

// GetUserInfo 获取用户信息
func (t *SimpleChaincode) GetUserInfo(ctx contractapi.TransactionContextInterface, userID string) (*User, error) {
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	if userBytes == nil {
		return nil, fmt.Errorf("user info not found")
	}
	user := &User{}
	if err = json.Unmarshal(userBytes, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %v", err)
	}
	return user, nil
}

// GetUsers 获取所有用户信息
func (t *SimpleChaincode) GetUsers(ctx contractapi.TransactionContextInterface) ([]*User, error) {
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey("user", []string{})
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	defer resultsIterator.Close()

	var users []*User
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get user info: %v", err)
		}
		user := &User{}
		if err = json.Unmarshal(response.Value, &user); err != nil {
			return nil, fmt.Errorf("failed to unmarshal user info: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

// UpdateUserInfo 更新用户信息
func (t *SimpleChaincode) UpdateUserInfo(ctx contractapi.TransactionContextInterface, userID string, username string, email string) error {
	userBytes, err := ctx.GetStub().GetState(userID)
	if err != nil {
		return fmt.Errorf("failed to get user info: %v", err)
	}
	if userBytes == nil {
		return fmt.Errorf("user info not found")
	}
	user := &User{}
	if err = json.Unmarshal(userBytes, &user); err != nil {
		return fmt.Errorf("failed to unmarshal user info: %v", err)
	}
	user.Username = username
	user.Email = email
	user.UpdatedAt = time.Now().Unix()
	userBytes, err = json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user info: %v", err)
	}
	if err = ctx.GetStub().PutState(user.ID, userBytes); err != nil {
		return fmt.Errorf("failed to PutState userBytes %s: %v", userBytes, err)
	}
	return nil
}
