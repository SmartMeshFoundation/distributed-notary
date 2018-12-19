package utils

import "io"
import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
)

//MyCallerFuncHandler handler for log
func MyCallerFuncHandler(h log.Handler) log.Handler {
	return log.FuncHandler(func(r *log.Record) error {
		r.Ctx = append(r.Ctx, "fn", fmt.Sprintf("%s:%n:%d", r.Call, r.Call, r.Call))
		return h.Log(r)
	})
}

//MyStreamHandler handler for log
func MyStreamHandler(wr io.Writer) log.Handler {
	fmtr := log.TerminalFormat(true)
	h := log.FuncHandler(func(r *log.Record) error {
		_, err := wr.Write(fmtr.Format(r))
		return err
	})
	return log.LazyHandler(log.SyncHandler(MyCallerFuncHandler(h)))
}
