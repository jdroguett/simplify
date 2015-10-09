package simplify

import (
	"log"
	"os"
)

type DebugLog struct {
	*log.Logger
}

var (
	Debug = false
	Log   = New()
)

func New() *DebugLog {
	dl := new(DebugLog)
	dl.Logger = log.New(os.Stderr, "[simplify] ", 1e9)
	return dl
}

func (dl *DebugLog) Println(v ...interface{}) {
	if Debug {
		dl.Logger.Println(v...)
	}
}
