package message

import (
	"github.com/ch4rl1e5/ccdb/internal/definitions"
	"github.com/ch4rl1e5/go-common/logger"
	"go.uber.org/zap"
	"io"
)

func ParseMessage(reader io.Reader) {
	buf := make([]byte, 4*1024)

	for {
		n, err := reader.Read(buf)
		if err != nil || n == 0 {
			logger.ZapLogger.Panic(definitions.ErrParseMessage, zap.Error(err))
		}

	}

}
