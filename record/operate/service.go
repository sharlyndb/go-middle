/**
 * @Time: 2022/3/7 15:36
 * @Author: yt.yin
 */

package operate

import (
	"github.com/golang-module/carbon/v2"
	"github.com/goworkeryyt/go-core/db"
	"github.com/goworkeryyt/go-core/global"
	"github.com/goworkeryyt/go-toolbox/page"
	"go.uber.org/zap"
)

type OperateRecordService struct{}

var OperateRecordServiceApp = new(OperateRecordService)

// CreateOperateRecord 创建记录
func (opt *OperateRecordService) CreateOperateRecord(record OperateRecord, retainDays int) (err error) {
	err = global.DB.Create(&record).Error
	go func() {
		// 默认保留100天
		if retainDays < 100 {
			retainDays = 100
		}
		time := carbon.Now().SubDays(retainDays).ToDateTimeString()
		err = global.DB.Where("create_time < ?", time).Delete(&OperateRecord{}).Error
		if err != nil {
			global.LOG.Error("删除操作记录异常：", zap.Any("err", err))
		}
	}()
	return err
}

// GetOperateRecordPage 分页获取操作记录列表
func (opt *OperateRecordService) GetOperateRecordPage(pageInfo *page.PageInfo) (err error, pageBean *page.PageBean) {
	pageBean = &page.PageBean{Page: pageInfo.Current, PageSize: pageInfo.RowCount}
	rows := make([]*OperateRecord, 0)
	err, pageBean = db.FindPage(&OperateRecord{}, &rows, pageInfo)
	return
}
