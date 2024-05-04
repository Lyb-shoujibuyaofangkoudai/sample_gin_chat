package middleware

import (
	"fmt"
	"gin_chat/global"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math"
	"os"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next() // 放行 调用该请求的剩余处理程序
		stopTime := time.Since(startTime)
		spendTime := fmt.Sprintf("%d ms", int(math.Ceil(float64(stopTime.Nanoseconds()/1000000))))
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "Unknown"
		}
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		dataSize := c.Writer.Size()
		if dataSize < 0 {
			dataSize = 0
		}
		method := c.Request.Method
		url := c.Request.RequestURI

		// 查询参数
		queryMap := c.Request.URL.Query()

		// 路径参数
		params := make(map[string]string)
		for _, param := range c.Params {
			params[param.Key] = param.Value
		}

		// 请求体
		body, _ := ioutil.ReadAll(c.Request.Body)
		requestBody := string(body)

		// 头信息
		headers := make(map[string][]string)
		for key, values := range c.Request.Header {
			headers[key] = values
		}

		// 将所有数据合并到一个map中
		allData := map[string]interface{}{
			"Query":   queryMap,
			"Params":  params,
			"Body":    requestBody,
			"Headers": headers,
			"Url":     url,
		}

		Log := global.Log.WithFields(logrus.Fields{
			"Type":      global.GormLog,
			"HostName":  hostName,
			"SpendTime": spendTime,
			"path":      url,
			"Method":    method,
			"status":    statusCode,
			"Ip":        clientIP,
			"DataSize":  dataSize,
			"UserAgent": userAgent,
			"AllData":   allData,
		})
		if len(c.Errors) > 0 {
			Log.Error(c.Errors.ByType(gin.ErrorTypePrivate))
		}
		msg := fmt.Sprintf("\n请求类型： %s \n请求地址： %s  \n请求query参数： %s \n请求body参数：%s", method, url, queryMap, requestBody)
		if statusCode >= 500 {
			Log.Error(msg)
		} else if statusCode >= 400 {
			Log.Warn(msg)
		} else {
			Log.Info(msg)
		}
	}
}
