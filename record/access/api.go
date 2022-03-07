/**
 * @Time: 2022/3/7 15:19
 * @Author: yt.yin
 */

package access

import (
	"github.com/gin-gonic/gin"
	"github.com/goworkeryyt/go-core/global"
	"github.com/goworkeryyt/go-toolbox/page"
	"github.com/goworkeryyt/go-toolbox/result"
	"go.uber.org/zap"
)

type AccessRecordApi struct{}

// GetAccessRecordPage 分页查询操作记录
func (s *AccessRecordApi) GetAccessRecordPage(c *gin.Context) {
	pageInfo := page.PageParam(c)
	if pageInfo == nil {
		result.FailMsg("获取失败,解析请求参数异常", c)
		return
	}
	err, pageBean := AccessRecordServiceApp.GetAccessRecordPage(pageInfo)
	if err != nil {
		global.LOG.Error("获取接口访问记录失败:", zap.Any("err", err))
		result.FailMsg("获取接口访问记录失败", c)
	} else {
		result.OkDataMsg(pageBean, "获取成功", c)
	}
}
