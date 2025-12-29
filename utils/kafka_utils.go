package utils

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"time"
)

func ParseAcks(acks string) int {
	switch acks {
	case "all":
		return -1
	case "1":
		return 1
	case "0":
		return 0
	default:
		logger.Info("invalid acks value '%s', using '1' as default", zap.String("acks", acks))
		return 1
	}
}

func FilterOffset(autoOffsetReset string) int64 {
	switch autoOffsetReset {
	case "earliest":
		return kafka.FirstOffset
	case "latest":
		return kafka.LastOffset
	default:
		return kafka.LastOffset
	}
}

func FilterEnableAutoCommit(enableAutoCommit bool) time.Duration {
	if enableAutoCommit {
		return 1 * time.Second
	}
	return 0
}

func FilterIsolationLevel(isolationLevel string) kafka.IsolationLevel {
	switch isolationLevel {
	case "read_uncommitted":
		return kafka.ReadUncommitted
	case "read_committed":
		return kafka.ReadCommitted
	default:
		return kafka.ReadCommitted
	}
}
