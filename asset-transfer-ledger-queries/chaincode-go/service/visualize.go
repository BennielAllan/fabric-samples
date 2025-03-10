package service

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

type TxHistory struct {
	TxID        string `json:"txID"`
	TxTimestamp string `json:"txTimestamp"`
	TxIsDelete  bool   `json:"txIsDelete"`
	TxValue     []byte `json:"txValue"`
}

func (t *SimpleChaincode) GetTxHistory(ctx contractapi.TransactionContextInterface, txID string) (*[]TxHistory, error) {
	history, err := ctx.GetStub().GetHistoryForKey(txID)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for key: %v", err)
	}
	defer history.Close()

	var txHistory []TxHistory
	for history.HasNext() {
		queryResponse, err := history.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get the next state for key %s: %v", txID, err)
		}
		txHistory = append(txHistory, TxHistory{
			TxID:        queryResponse.TxId,
			TxTimestamp: queryResponse.Timestamp.String(),
			TxIsDelete:  queryResponse.IsDelete,
			TxValue:     queryResponse.Value,
		})
	}
	return &txHistory, nil
}
