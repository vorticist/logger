package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
}

func SetOutputFile(writer io.Writer) {
	log.SetOutput(writer)
}

func Trace(args ...interface{}) {
	log.Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	logf(format, log.Tracef, args)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logf(format, log.Debugf, args)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logf(format, log.Infof, args)
}

func Print(args ...interface{}) {
	log.Print(args...)
}

func Printf(format string, args ...interface{}) {
	logf(format, log.Printf, args)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logf(format, log.Warnf, args)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logf(format, log.Errorf, args)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	logf(format, log.Fatalf, args)
}

func Panic(args ...interface{}) {
	log.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	logf(format, log.Panicf, args)
}

func logf(format string, logf func(format string, args ...interface{}), args []interface{}) {
	var jsonArgs []interface{}

	for _, arg := range args {
		if req, ok := arg.(*http.Request); ok {
			reqData := extractRequestData(req)
			jsonArg, err := json.Marshal(reqData)
			if err != nil {
				jsonArgs = append(jsonArgs, fmt.Sprintf("error marshaling http.Request: %v", err))
			} else {
				jsonArgs = append(jsonArgs, string(jsonArg))
			}
			continue
		}

		if reflect.TypeOf(arg).Kind() == reflect.Func {
			signature := getFunctionSignature(arg)
			jsonArgs = append(jsonArgs, signature)
		} else {
			value := dereferencePointer(arg)

			jsonArg, err := json.Marshal(value)
			if err != nil {
				jsonArgs = append(jsonArgs, fmt.Sprintf("error marshaling arg: %v", err))
			} else {
				jsonArgs = append(jsonArgs, string(jsonArg))
			}
		}
	}

	logf(format, jsonArgs...)
}
func getFunctionSignature(fn interface{}) string {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return "not a function"
	}

	var params []string
	for i := 0; i < fnType.NumIn(); i++ {
		params = append(params, fnType.In(i).String())
	}

	var returns []string
	for i := 0; i < fnType.NumOut(); i++ {
		returns = append(returns, fnType.Out(i).String())
	}

	return fmt.Sprintf("func(%s) (%s)", strings.Join(params, ", "), strings.Join(returns, ", "))
}

func dereferencePointer(arg interface{}) interface{} {
	val := reflect.ValueOf(arg)

	if val.Kind() == reflect.Ptr {
		if !val.IsNil() {
			return val.Elem().Interface()
		}
		return "nil pointer"
	}

	return arg
}

func extractRequestData(req *http.Request) map[string]interface{} {
	// Read the body (if it's not already read)
	var bodyData string
	if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err == nil {
			bodyData = string(bodyBytes)
		} else {
			bodyData = "error reading body"
		}
		// Restore the io.ReadCloser by re-creating the body
		req.Body = io.NopCloser(strings.NewReader(bodyData))
	}

	// Return a map with the relevant request data
	return map[string]interface{}{
		"method":  req.Method,
		"url":     req.URL.String(),
		"headers": req.Header,
		"body":    bodyData,
	}
}
