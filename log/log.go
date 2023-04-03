package log

import (
	"gin-demo/core/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net"
)

// Logger logger
type Logger = *logrus.Entry

var Loger Logger

var ServerIP = localIP()

func init() {
	Loger = logrus.WithFields(logrus.Fields{
		"app_name":    config.AppName,
		"server_host": config.AppAddr,
		"env":         config.Env,
		"server_ip":   ServerIP,
	})

	Loger.Logger.SetReportCaller(true)
}

func init() {
	setLogConfig()

	//文件日志Hook
	NewFileHook()

	//kafka Hook
	//if config.IsSet("log.kafka") {
	//	kafka := config.GetStringSlice("log.kafka")
	//
	//	NewKafkaHook(config.AppName, &logrus.JSONFormatter{}, kafka)
	//}
}

func setLogConfig() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	if config.LogFormatter == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	//设置日志输出级别
	logLevel := getLogLevel(config.LogLevel)
	logrus.SetLevel(logLevel)
}

func getLogLevel(lvl string) logrus.Level {
	level, ok := logrus.ParseLevel(lvl)
	if ok == nil {
		return level
	}

	if config.IsDevEnv { // 开发环境设置为debug级别
		return logrus.DebugLevel
	}

	return logrus.ErrorLevel
}

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// Get 获取日志实例
func Get(c *gin.Context) Logger {
	//获取request信息
	if c == nil || c.Request == nil {
		return Loger.WithFields(nil)
	}

	body := ""
	// 获取 gin.Context 缓存的request 请求参数。
	// 如果未缓存也不从 c.Request 中获取，避免携程并发读取冲突问题
	if b, ok := c.Get(gin.BodyBytesKey); ok {
		body = b.(string)
	}

	return logrus.WithFields(logrus.Fields{
		"request_id":    c.GetString("X-Request-ID"),
		"app_ame":       config.AppName,
		"domain":        c.Request.Host,
		"referer":       c.Request.Referer(),
		"client_ip":     c.ClientIP(),
		"method":        c.Request.Method,
		"url":           c.Request.RequestURI,
		"user_agent":    c.Request.UserAgent(),
		"request_param": body,
	})
}
