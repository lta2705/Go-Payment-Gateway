package handler

import (
	"encoding/json"
	"github.com/Jeffail/gabs"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/lta2705/Go-Payment-Gateway/internal/constant"
	"github.com/lta2705/Go-Payment-Gateway/internal/dto"
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

	msgType, ok := jsonParsed.Path("msgType").Data().(string)
	if !ok {
		return nil
	}

	logger.Info("Message type:", msgType)

	switch msgType {
	case constant.MsgTypeTxReq:
		logger.Info("Processing transaction request message")
		return h.updateTransaction(payload, func(txDto *dto.TerminalTransactionDTO) {
			txDto.UpdatedBy = "TCP_SERVER"
			txDto.ErrorCode = constant.ErrCodeTrmNotResponse
			txDto.ErrorDetail = constant.ErrDetailCode11
			txDto.Status = constant.TxStatusFailed

		})
	case constant.MsgTypeTxRes:
		logger.Info("Processing transaction response message")
		return h.updateTransaction(payload, func(txDto *dto.TerminalTransactionDTO) {
			txDto.UpdatedBy = "TERMINAL"
			txDto.ErrorCode = constant.ErrCodeNoErr
			txDto.ErrorDetail = constant.ErrDetailCode0
		})
	}
	return nil
}

func (h *TransactionRespHandlerImpl) updateTransaction(payload []byte, updateFn func(*dto.TerminalTransactionDTO)) error {
	var txDto dto.TerminalTransactionDTO
	if err := json.Unmarshal(payload, &txDto); err != nil {
		logger.Error("Error when parsing payload", err)
		return err
	}

	if updateFn != nil {
		updateFn(&txDto)
	}

	tx, err := h.repository.FindByPcPosIdAndTransactionId(
		txDto.PcPosId,
		txDto.TransactionId,
	)
	if err != nil {
		logger.Error("Error when finding transaction:", err)
		return err
	}

	if tx == nil {
		logger.Warn("Transaction not found in database for PcPosId and TransactionId")
		return nil
	}

	logger.Info("Found transaction:", tx)
	logger.Info("transaction dto:", txDto)

	tx.UpdatedBy = txDto.UpdatedBy
	tx.Status = txDto.Status
	tx.ErrorCode = txDto.ErrorCode
	tx.ErrorDetail = txDto.ErrorDetail

	logger.Info("Updating transaction:", tx)
	return h.repository.UpdateTransaction(tx)
}
