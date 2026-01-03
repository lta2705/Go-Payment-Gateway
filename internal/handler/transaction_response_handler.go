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
	logger.Info("Handling the transaction response:", string(payload))
	jsonParsed, err := gabs.ParseJSON(payload)
	if err != nil {
		return err
	}

	logger.Info("Parsed JSON:", jsonParsed.String())

	msgType, ok := jsonParsed.Path("MsgType").Data().(string)
	if !ok {
		return nil
	}

	logger.Info("Message type:", msgType)

	switch msgType {
	case constant.MsgTypeTxReq:
		return h.updateTransaction(payload, func(tx *dto.TerminalTransactionDTO) {
			tx.UpdatedBy = "TCP_SERVER"
			tx.ErrorCode = constant.ErrCodeTrmNotResponse
			tx.ErrorDetail = constant.ErrDetailCode11
			tx.Status = constant.TxStatusFailed

		})
	case constant.MsgTypeTxRes:
		return h.updateTransaction(payload, func(tx *dto.TerminalTransactionDTO) {
			tx.UpdatedBy = "TERMINAL"
			tx.ErrorCode = constant.ErrCodeNoErr
			tx.ErrorDetail = constant.ErrDetailCode0
			tx.Status = constant.TxStatusSuccess
		})
	}
	return nil
}

func (h *TransactionRespHandlerImpl) updateTransaction(payload []byte, updateFn func(*dto.TerminalTransactionDTO)) error {
	var txDto dto.TerminalTransactionDTO
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

	if tx.ID == [16]byte{} && txDto.ID != [16]byte{} {
		tx.ID = txDto.ID
	}

	if tx.PcPosId == "" && txDto.PcPosId != "" {
		tx.PcPosId = txDto.PcPosId
	}

	logger.Info("Updating transaction:", tx.ToBeautifiedString())
	return h.repository.UpdateTransaction(&tx)
}
