package repository

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/lta2705/Go-Payment-Gateway/internal/model"
	"gorm.io/gorm"
)

type MerchantTerminalRepository interface {
	FindTerminalIdByPcPosId(pcPosId string) (string, error)
}

type MerchantTerminalRepositoryImpl struct {
	db *gorm.DB
}

func (m *MerchantTerminalRepositoryImpl) FindTerminalIdByPcPosId(pcPosId string) (string, error) {
	var terminal model.MerchantTerminal

	err := m.db.Select("terminal_id").Where("pc_pos_id = ?", pcPosId).First(&terminal).Error
	if err != nil {
		logger.Warn("Cannot find terminalID by pcPosId", pcPosId)
		return "", err
	}

	return terminal.TerminalId, nil
}
