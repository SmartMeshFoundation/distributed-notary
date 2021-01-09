package log

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
)

var (
	root = &logger{[]interface{}{}, new(swapHandler)}
	//StdoutHandler handler output to stdout
	StdoutHandler = StreamHandler(os.Stdout, LogfmtFormat())
	//StderrHandler output to stderr
	StderrHandler = StreamHandler(os.Stderr, LogfmtFormat())
)

func init() {
	//root.SetHandler(DiscardHandler())
	Root().SetHandler(LvlFilterHandler(LvlTrace, DefaultStreamHandler(os.Stdout)))
}

// New returns a new logger with the given context.
// New is a convenient alias for Root().New
func New(ctx ...interface{}) Logger {
	return root.New(ctx...)
}

// Root returns the root logger
func Root() Logger {
	return root
}

// The following functions bypass the exported logger methods (logger.Debug,
// etc.) to keep the call depth the same for all paths to logger.write so
// runtime.Caller(2) always refers to the call site in client code.

// Trace is a convenient alias for Root().Trace
func Trace(format string, ctx ...interface{}) {
	root.write(format, LvlTrace, ctx )
}
func Tracef(format string, args ...interface{}) {
	format = fmt.Sprintf(format, args...)
	root.write(format, LvlTrace, nil)
}
// Debug is a convenient alias for Root().Debug
func Debug(format string, ctx ...interface{}) {
	root.write(format, LvlDebug, ctx)
}
// Debugf is a convenient alias for Root().Debug
func Debugf(format string, args ...interface{}) {
	format = fmt.Sprintf(format, args...)
	root.write(format, LvlDebug, nil)
}
// Info is a convenient alias for Root().Info
func Info(format string, ctx ...interface{}) {
	root.write(format, LvlInfo, ctx)
}
// Infof is a convenient alias for Root().Info
func Infof(format string, args ...interface{}) {
	format = fmt.Sprintf(format, args...)
	root.write(format, LvlInfo, nil)
}
// Warn is a convenient alias for Root().Warn
func Warn(format string, ctx ...interface{}) {
	root.write(format, LvlWarn, ctx)
}
// Warnf is a convenient alias for Root().Warn
func Warnf(format string, args ...interface{}) {
	format = fmt.Sprintf(format, args...)
	root.write(format, LvlWarn, nil)
}

// Error is a convenient alias for Root().Error
func Error(format string, ctx ...interface{}) {
	root.write(format, LvlError, ctx)
}
// Errorf is a convenient alias for Root().Error
func Errorf(format string, args ...interface{}) {
	format = fmt.Sprintf(format, args...)
	root.write(format, LvlError, nil)
}
// Crit is a convenient alias for Root().Crit
func Crit(format string, ctx ...interface{}) {
	root.write(format, LvlCrit, ctx)
	os.Exit(1)
}
// Crit is a convenient alias for Root().Crit
func Critf(format string, args ...interface{}) {
	format = fmt.Sprintf(format, args...)
	root.write(format, LvlCrit, nil)
	os.Exit(1)
}
//StringInterface conver any object to string and it's depth is `depth`
func StringInterface(i interface{}, depth int) string {
	stringer, ok := i.(fmt.Stringer)
	if ok {
		return stringer.String()
	}
	//c := spew.Config
	//spew.Config.DisableMethods = false
	////spew.Config.ContinueOnMethod = false
	spew.Config.MaxDepth = depth
	s := spew.Sdump(i)
	return s
}

//StringInterface1 convert any object to string ,it's depth is 1
func StringInterface1(i interface{}) string {
	stringer, ok := i.(fmt.Stringer)
	if ok {
		return stringer.String()
	}
	//c := spew.Config
	//spew.Config.DisableMethods = true
	spew.Config.MaxDepth = 1
	s := spew.Sdump(i)
	//	spew.Config = c
	return s
}
