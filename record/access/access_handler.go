/**
 * @Time: 2022/3/7 14:57
 * @Author: yt.yin
 */

package access

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"github.com/goworkeryyt/go-core/global"
	"github.com/goworkeryyt/go-toolbox/uuid"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// AccessRecordHandler 访问记录 retainDays 请求记录保留的时间
func AccessRecordHandler(retainDays int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		httpMethod := c.Request.Method
		if httpMethod != http.MethodGet {
			var err error
			body, err = ioutil.ReadAll(c.Request.Body)
			if err != nil {
				global.LOG.Error("read body from request error:", zap.Any("err", err))
			} else {
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}
		}
		header := c.Request.Header.Get("Content-Type")
		record := AccessRecord{
			Ip:     c.ClientIP(),
			Method: httpMethod,
			Path:   c.Request.URL.Path,
			Agent:  c.Request.UserAgent(),
			Body:   string(body),
		}
		if strings.Contains(header, "multipart/form-data") {
			// 如果该请求是文件上传不记录请求体
			record.Body = "文件上传"
		} else {
			record.Body = string(body)
		}
		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		now := time.Now()
		c.Next()
		record.Error = c.Errors.ByType(gin.ErrorTypePrivate).String()
		record.Status = c.Writer.Status()
		record.Latency = time.Now().Sub(now).Milliseconds()
		record.Resp = writer.body.String()
		record.ID = uuid.UUID()
		record.CreateTime = carbon.Now().ToDateTimeString()
		if httpMethod == http.MethodGet {
			var str bytes.Buffer
			m, _ := json.Marshal(record)
			_ = json.Indent(&str, m, "", "    ")
			global.LOG.Info("访问记录" + str.String())
		}
		if err := AccessRecordServiceApp.CreateAccessRecord(record, retainDays); err != nil {
			global.LOG.Error("create access record error:", zap.Any("err", err))
		}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
