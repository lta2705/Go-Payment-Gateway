package handler

import (
	"encoding/json"
	"github.com/Jeffail/gabs"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/jinzhu/copier"
	"github.com/lta2705/Go-Payment-Gateway/internal/constant"
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"github.com/lta2705/Go-Payment-Gateway/internal/repository"
)

type TransactionRespHandler interface {
	HandleTransaction(payload []byte) error
}

type TransactionRespHandlerImpl struct {
	repository repository.TransactionRepository
}

func NewTransactionRespHandler(repo repository.TransactionRepository) TransactionRespHandler {
	return &TransactionRespHandlerImpl{repository: repo}
}

func (h *TransactionRespHandlerImpl) HandleTransaction(payload []byte) error {
	jsonParsed, err := gabs.ParseJSON(payload)
	if err != nil {
		return err
	}

	msgType, ok := jsonParsed.Path("msgType").Data().(string)
	if !ok {
		return nil
	}

	switch msgType {
	case constant.MsgTypeTxReq:
		return h.updateTransaction(payload, func(tx *dto.TransactionDTO) {
			tx.UpdatedBy = "TCP_SERVER"
			tx.ErrorCode = constant.ErrCodeTrmNotResponse
			tx.ErrorDetail = constant.ErrDetailCode11
			tx.Status = constant.TxStatusFailed

		})
	case constant.MsgTypeTxRes:
		return h.updateTransaction(payload, func(tx *dto.TransactionDTO) {
			tx.UpdatedBy = "TERMINAL"
			tx.ErrorCode = constant.ErrCodeNoErr
			tx.ErrorDetail = constant.ErrDetailCode0
			tx.Status = constant.TxStatusSuccess
		})
	}
	return nil
}

// Giữ hàm updateTransaction nội bộ trong handler
func (h *TransactionRespHandlerImpl) updateTransaction(payload []byte, updateFn func(*dto.TransactionDTO)) error {
	var txDto dto.TransactionDTO
	err := json.Unmarshal(payload, &txDto)
	if err != nil {
		logger.Error("Error when parsing payload", err)
		return err
	}
	updateFn(&txDto)

	var tx model.Transaction
	copyErr := copier.Copy(&tx, &txDto)
	if copyErr != nil {
		logger.Error("Error when copying struct:", copyErr)
		return copyErr
	}
	return h.repository.UpdateTransaction(&tx)
}
