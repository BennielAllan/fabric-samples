package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"time"
)

const (
	StatusPending   = "pending"
	StatusActive    = "active"
	StatusCompleted = "completed"
	StatusExpired   = "expired"
)

// RaiseProject 筹款项目表
type RaiseProject struct {
	DocType       string  `json:"docType"`
	ID            string  `json:"id"`
	BeneficiaryID string  `json:"beneficiary_id"` // 筹款人id
	TargetAmount  float64 `json:"target_amount"`  // 目标金额
	CurrentAmount float64 `json:"current_amount"` // 已筹款金额
	MedicalProof  string  `json:"medical_proof"`  // 医疗证明
	Status        string  `json:"status"`         // pending（待审核）、active（进行中）、completed（已完成）、expired（已过期）
	CreatedAt     int64   `json:"created_at"`
	UpdatedAt     int64   `json:"updated_at"`
	Deadline      int64   `json:"deadline"`
	MedicalBillID int64   `json:"medical_bill_id"` // 关联的账单id
}

// MedicalBills 医疗账单表
type MedicalBills struct {
	DocType       string  `json:"docType"`
	BillID        string  `json:"bill_id"`        // 账单ID
	PatientID     string  `json:"campaign_id"`    // 患者ID
	InstitutionID string  `json:"institution_id"` // 医疗机构ID
	Amount        float64 `json:"amount"`         // 金额
	FilePath      string  `json:"file_path"`      // 证明文件路径
	Status        string  `json:"status"`         // pending（待审核）、active（进行中）、completed（已完成）、expired（已过期）
	CreatedAt     int64   `json:"created_at"`
	UpdatedAt     int64   `json:"updated_at"`
}

// Transaction 交易记录表
type Transaction struct {
	DocType        string  `json:"docType"`
	ID             string  `json:"id"`               // 交易id
	RaiseProjectID string  `json:"raise_project_id"` // 筹款项目id
	DonorID        string  `json:"donor_id"`         // 捐赠者id
	Amount         float64 `json:"amount"`           // 金额
	Type           string  `json:"type"`             // 类型
	BlockchainHash string  `json:"blockchain_hash"`  // 区块链hash
	Timestamp      int64   `json:"timestamp"`
}

// CreateRaiseProject 创建筹款项目
func (t *SimpleChaincode) CreateRaiseProject(ctx contractapi.TransactionContextInterface,
	beneficiaryID string, targetAmount float64, medicalProof string, deadline int64) error {
	raiseProject := &RaiseProject{
		DocType:       "raise_project",
		BeneficiaryID: beneficiaryID,
		TargetAmount:  targetAmount,
		MedicalProof:  medicalProof,
		Status:        StatusPending,
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
		Deadline:      deadline,
	}
	raiseProjectBytes, err := json.Marshal(raiseProject)
	if err != nil {
		return fmt.Errorf("failed to marshal raise project info: %v", err)
	}
	if err = ctx.GetStub().PutState(raiseProject.ID, raiseProjectBytes); err != nil {
		return fmt.Errorf("failed to PutState raiseProjectBytes %s: %v", raiseProjectBytes, err)
	}
	return nil
}

// CreateMedicalBill 创建医疗账单
func (t *SimpleChaincode) CreateMedicalBill(ctx contractapi.TransactionContextInterface,
	patientID string, institutionID string, amount float64, filePath string) error {
	medicalBills := &MedicalBills{
		BillID:        uuid.NewString(),
		PatientID:     patientID,
		InstitutionID: institutionID,
		Amount:        amount,
		FilePath:      filePath,
		Status:        StatusPending,
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     time.Now().Unix(),
	}
	medicalBillsBytes, err := json.Marshal(medicalBills)
	if err != nil {
		return fmt.Errorf("failed to marshal medical bills info: %v", err)
	}
	if err = ctx.GetStub().PutState(medicalBills.BillID, medicalBillsBytes); err != nil {
		return fmt.Errorf("failed to PutState medicalBillsBytes %s: %v", medicalBillsBytes, err)
	}
	return nil
}

// GetRaiseProjectByPId 根据筹款项目ID获取筹款项目
func (t *SimpleChaincode) GetRaiseProjectByPId(ctx contractapi.TransactionContextInterface, raiseProjectID string) (*RaiseProject, error) {
	raiseProjectBytes, err := ctx.GetStub().GetState(raiseProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to GetState raiseProjectBytes %s: %v", raiseProjectID, err)
	}
	if raiseProjectBytes == nil {
		return nil, fmt.Errorf("raiseProjectBytes %s does not exist", raiseProjectID)
	}
	var raiseProject RaiseProject
	if err := json.Unmarshal(raiseProjectBytes, &raiseProject); err != nil {
		return nil, fmt.Errorf("failed to unmarshal raiseProjectBytes %s: %v", raiseProjectBytes, err)
	}
	return &raiseProject, nil
}

// GetRaiseProjectByUId 根据用户ID获取所有筹款项目
func (t *SimpleChaincode) GetRaiseProjectByUId(ctx contractapi.TransactionContextInterface, userID string) ([]*RaiseProject, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"raise_project","beneficiary_id":"%s"}}`, userID)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to GetQueryResult for query %s: %v", queryString, err)
	}
	defer resultsIterator.Close()

	var raiseProjects []*RaiseProject
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to query next result: %v", err)
		}
		var raiseProject RaiseProject
		if err := json.Unmarshal(queryResponse.Value, &raiseProject); err != nil {
			return nil, fmt.Errorf("failed to unmarshal raiseProjectBytes %s: %v", queryResponse.Value, err)
		}
		raiseProjects = append(raiseProjects, &raiseProject)
	}
	return raiseProjects, nil
}

// GetRaiseProjects 获取所有筹款项目
func (t *SimpleChaincode) GetRaiseProjects(ctx contractapi.TransactionContextInterface) ([]*RaiseProject, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"raise_project"}}`)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to GetQueryResult for query %s: %v", queryString, err)
	}
	defer resultsIterator.Close()

	var raiseProjects []*RaiseProject
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to query next result: %v", err)
		}
		var raiseProject RaiseProject
		if err := json.Unmarshal(queryResponse.Value, &raiseProject); err != nil {
			return nil, fmt.Errorf("failed to unmarshal raiseProjectBytes %s: %v", queryResponse.Value, err)
		}
		raiseProjects = append(raiseProjects, &raiseProject)
	}
	return raiseProjects, nil
}

// Recharge 充值
func (t *SimpleChaincode) Recharge(ctx contractapi.TransactionContextInterface, userID string, amount float64) error {
	user, err := t.GetUserInfo(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to GetUserInfo %s: %v", userID, err)
	}
	user.Balance += amount
	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user info: %v", err)
	}
	if err = ctx.GetStub().PutState(user.ID, userBytes); err != nil {
		return fmt.Errorf("failed to PutState userBytes %s: %v", userBytes, err)
	}
	return nil
}

// Donate 捐助
func (t *SimpleChaincode) Donate(ctx contractapi.TransactionContextInterface, donorID, raiseProjectID string, amount float64) error {
	// 更新用户余额
	user, err := t.GetUserInfo(ctx, donorID)
	if err != nil {
		return fmt.Errorf("failed to GetUserInfo %s: %v", donorID, err)
	}
	if user.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}
	user.Balance -= amount
	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user info: %v", err)
	}
	if err = ctx.GetStub().PutState(user.ID, userBytes); err != nil {
		return fmt.Errorf("failed to PutState userBytes %s: %v", userBytes, err)
	}

	// 更新筹款项目
	raiseProject, err := t.GetRaiseProjectByPId(ctx, raiseProjectID)
	if err != nil {
		return fmt.Errorf("failed to GetRaiseProjectByPId %s: %v", raiseProjectID, err)
	}
	raiseProject.CurrentAmount += amount
	if raiseProject.CurrentAmount >= raiseProject.TargetAmount {
		raiseProject.Status = StatusCompleted
	}
	raiseProjectBytes, err := json.Marshal(raiseProject)
	if err != nil {
		return fmt.Errorf("failed to marshal raise project info: %v", err)
	}
	if err = ctx.GetStub().PutState(raiseProject.ID, raiseProjectBytes); err != nil {
		return fmt.Errorf("failed to PutState raiseProjectBytes %s: %v", raiseProjectBytes, err)
	}

	// 记录交易
	transaction := Transaction{
		DocType:        "transaction",
		ID:             uuid.NewString(),
		Amount:         amount,
		RaiseProjectID: raiseProjectID,
		DonorID:        donorID,
		Timestamp:      time.Now().Unix(),
	}
	transactionBytes, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal raise project info: %v", err)
	}
	if err = ctx.GetStub().PutState(transaction.ID, transactionBytes); err != nil {
		return fmt.Errorf("failed to PutState raiseProjectBytes %s: %v", transactionBytes, err)
	}
	return nil
}

// GetTxByUid 根据用户id获取交易列表
func (t *SimpleChaincode) GetTxByUid(ctx contractapi.TransactionContextInterface, userID string) ([]*Transaction, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"transaction","donor_id":"%s"}}`, userID)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to GetQueryResult for query %s: %v", queryString, err)
	}
	defer resultsIterator.Close()

	var transactions []*Transaction
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to query next result: %v", err)
		}
		var transaction Transaction
		if err := json.Unmarshal(queryResponse.Value, &transaction); err != nil {
			return nil, fmt.Errorf("failed to unmarshal raiseProjectBytes %s: %v", queryResponse.Value, err)
		}
		transactions = append(transactions, &transaction)
	}
	return transactions, nil
}

// GetTxByPId 根据筹款项目id获取交易列表
func (t *SimpleChaincode) GetTxByPId(ctx contractapi.TransactionContextInterface, raiseProjectID string) ([]*Transaction, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"transaction","raise_project_id":"%s"}}`, raiseProjectID)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to GetQueryResult for query %s: %v", queryString, err)
	}
	defer resultsIterator.Close()

	var transactions []*Transaction
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to query next result: %v", err)
		}
		var transaction Transaction
		if err := json.Unmarshal(queryResponse.Value, &transaction); err != nil {
			return nil, fmt.Errorf("failed to unmarshal raiseProjectBytes %s: %v", queryResponse.Value, err)
		}
		transactions = append(transactions, &transaction)
	}
	return transactions, nil
}
