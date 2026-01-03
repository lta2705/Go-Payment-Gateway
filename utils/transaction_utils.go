package utils

import (
	"encoding/json"
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
)

func EnrichTransactionToJSON(txn *model.Transaction, terminalId string) (string, error) {
	data, err := json.Marshal(txn)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}

	result["TerminalId"] = terminalId

	finalData, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(finalData), nil
}
