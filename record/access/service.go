/**
 * @Time: 2022/3/7 15:08
 * @Author: yt.yin
 */

package access

import (
	"github.com/golang-module/carbon/v2"
	"github.com/goworkeryyt/go-core/db"
	"github.com/goworkeryyt/go-core/global"
	"github.com/goworkeryyt/go-toolbox/page"
	"go.uber.org/zap"
)

type AccessRecordService struct{}

var AccessRecordServiceApp = new(AccessRecordService)

// CreateAccessRecord 创建记录
func (opt *AccessRecordService) CreateAccessRecord(record AccessRecord, retainDays int) (err error) {
	err = global.DB.Create(&record).Error
	go func() {
		// 默认保留7天
		if retainDays < 7 {
			retainDays = 7
		}
		time := carbon.Now().SubDays(retainDays).ToDateTimeString()
		err = global.DB.Where("create_time < ?", time).Delete(&AccessRecord{}).Error
		if err != nil {
			global.LOG.Error("删除访问记录异常：", zap.Any("err", err))
		}
	}()
	return err
}

// GetAccessRecordPage 分页获取操作记录列表
func (opt *AccessRecordService) GetAccessRecordPage(pageInfo *page.PageInfo) (err error, pageBean *page.PageBean) {
	pageBean = &page.PageBean{Page: pageInfo.Current, PageSize: pageInfo.RowCount}
	rows := make([]*AccessRecord, 0)
	err, pageBean = db.FindPage(&AccessRecord{}, &rows, pageInfo)
	return
}
