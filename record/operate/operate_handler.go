/**
 * @Time: 2022/3/7 15:30
 * @Author: yt.yin
 */

package operate

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"github.com/goworkeryyt/go-core/global"
	"github.com/goworkeryyt/go-core/jwt"
	"github.com/goworkeryyt/go-toolbox/uuid"
	"go.uber.org/zap"
)

func OperateRecordHandler(retainDays int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		var userId string
		var username string
		var merchantNo string
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
		if claims, ok := c.Get("claims"); ok {
			waitUse := claims.(*jwt.CustomClaims)
			userId = waitUse.UserId
			username = waitUse.Username
			merchantNo = waitUse.MerchantNo
		} else {
			token := c.Request.Header.Get("ACCESS_TOKEN")
			if token != "" {
				j := jwt.NewJWT()
				// 解析token包含的信息
				claims, err := j.ResolveToken(token)
				if err == nil {
					userId = claims.UserId
					username = claims.Username
					merchantNo = claims.MerchantNo
				}
			}
		}
		record := OperateRecord{
			Ip:         c.ClientIP(),
			Method:     httpMethod,
			Path:       c.Request.URL.Path,
			Agent:      c.Request.UserAgent(),
			UserID:     userId,
			MerchantNo: merchantNo,
			Username:   username,
		}
		if strings.Contains(header, "multipart/form-data") {
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
		if err := OperateRecordServiceApp.CreateOperateRecord(record,retainDays); err != nil {
			global.LOG.Error("create operate record error:", zap.Any("err", err))
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
