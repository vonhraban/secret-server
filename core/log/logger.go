package log

type Logger interface {
	Info(args ...interface{}) 
	Debug(args ...interface{}) 
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})


	Infof(format string, args ...interface{}) 
	Debugf(format string, args ...interface{}) 
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{}) 
	Fatalf(format string, args ...interface{}) 
}

