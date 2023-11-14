package logs

import (
	"go.uber.org/zap"
)

func InitDefaultLogger() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}
