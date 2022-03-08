/**
 * @Time: 2022/3/7 14:56
 * @Author: yt.yin
 */

package record

import (
	"github.com/gin-gonic/gin"
	"github.com/goworkeryyt/go-core/global"
	"github.com/goworkeryyt/go-middle/record/access"
	"github.com/goworkeryyt/go-middle/record/operate"
	"go.uber.org/zap"
)

type ApiGroup struct {
	access.AccessRecordApi
	operate.OperateRecordApi
}

var ApiGroupApp = new(ApiGroup)

// autoCreateTables 初始化的时候自定创建表
func autoCreateTables() {
	if global.DB != nil {
		// 数据库自动迁移
		err := global.DB.AutoMigrate(
			operate.OperateRecord{},
			access.AccessRecord{},
		)
		if err != nil && global.LOG != nil {
			global.LOG.Error("初始化表时异常：", zap.Any("err", err))
		}
	}
}

// RouterRegister 注册访问记录中间件路由
func RouterRegister(rGroup *gin.RouterGroup) {
	// 初始化表
	autoCreateTables()
	// 创建路由
	{
		accessRecordApi := ApiGroupApp.AccessRecordApi
		operateRecodeApi := ApiGroupApp.OperateRecordApi
		// 私有接口
		recordGroup := rGroup.Group("record")
		{
			recordGroup.GET("getAccessRecordPage", accessRecordApi.GetAccessRecordPage)
			recordGroup.GET("getOperateRecordPage", operateRecodeApi.GetOperateRecordPage)
		}
	}
}
