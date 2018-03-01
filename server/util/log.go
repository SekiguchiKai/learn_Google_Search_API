package util

import (
	"strconv"
	"context"
	"runtime"
	"google.golang.org/appengine/log"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
)

func RespondAndLog(c *gin.Context, code int, format string, values ...interface{}) {
	if code >= 500 {
		ErrorLog(c, format, values...)
	} else if code >= 400 {
		InfoLog(c, format, values...)
	}
	c.String(code, format, values...)
}


func CriticalLog(c *gin.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "nofile"
		line = -1
	}
	ctx := getGoContextFromGinContext(c)
	log.Criticalf(ctx, file+":"+strconv.Itoa(line)+":"+format, args...)
}

func DebugLog(c *gin.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "nofile"
		line = -1
	}
	ctx := getGoContextFromGinContext(c)
	log.Debugf(ctx, file+":"+strconv.Itoa(line)+":"+format, args...)
}

func ErrorLog(c *gin.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "nofile"
		line = -1
	}
	ctx := getGoContextFromGinContext(c)
	log.Errorf(ctx, file+":"+strconv.Itoa(line)+":"+format, args...)
}

func InfoLog(c *gin.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "nofile"
		line = -1
	}

	ctx := getGoContextFromGinContext(c)
	log.Infof(ctx, file+":"+strconv.Itoa(line)+":"+format, args...)
}

func WarningLog(c *gin.Context, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "nofile"
		line = -1
	}

	ctx := getGoContextFromGinContext(c)
	log.Warningf(ctx, file+":"+strconv.Itoa(line)+":"+format, args...)
}

func getGoContextFromGinContext(c *gin.Context)context.Context {
	r := c.Request
	return appengine.NewContext(r)
}