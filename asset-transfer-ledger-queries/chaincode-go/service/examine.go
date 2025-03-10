package service

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"time"
)

const (
	TypeRaiseProject = "raise_project"
	TypeInstitution  = "institution"
)

// ExamineLogs 审核记录表
type ExamineLogs struct {
	ID         string `json:"id"`
	TargetType string `json:"target_type"` // raise_project（求助项目）、institution（医疗机构）、medical_bill（医疗账单）
	TargetID   string `json:"target_id"`   // 目标id
	AuditorID  string `json:"auditor_id"`  // 审核者ID
	Decision   bool   `json:"decision"`    // true为通过，false为不通过
	Comments   string `json:"comments"`    // 审核意见
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

// Examine 审核筹款项目或者医疗机构
func (t *SimpleChaincode) Examine(ctx contractapi.TransactionContextInterface, targetID, targetType, comments string, decision bool) error {
	userID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get user ID: %v", err)
	}

	examineLog := ExamineLogs{
		TargetType: targetType,
		TargetID:   targetID,
		AuditorID:  userID,
		Decision:   decision,
		Comments:   comments,
		CreatedAt:  time.Now().Unix(),
	}
	examineLogByte, err := json.Marshal(examineLog)
	if err != nil {
		return fmt.Errorf("failed to marshal examineLog: %v", err)
	}
	if err = ctx.GetStub().PutState(examineLog.ID, examineLogByte); err != nil {
		return fmt.Errorf("failed to PutState examineLog %s: %v", examineLogByte, err)
	}
	return nil
}
