/**
 * @Time: 2022/3/7 14:57
 * @Author: yt.yin
 */

package access

type AccessRecord struct {
	ID string `json:"id"             gorm:"column:id; primary_key;type:varchar(36)"`

	/** 创建时间 */
	CreateTime string `json:"createTime"     gorm:"column:create_time;index;type:varchar(20)"`

	/** 请求ip */
	Ip string `json:"ip"             gorm:"column:ip;comment:请求ip"`

	/** 请求方法 */
	Method string `json:"method"         gorm:"column:method;comment:请求方法"`

	/** 请求路径 */
	Path string `json:"path"           gorm:"column:path;comment:请求路径"`

	/** 请求状态 */
	Status int `json:"status"         gorm:"column:status;comment:请求状态"`

	/** 延迟 */
	Latency int64 `json:"latency"        gorm:"column:latency;comment:延迟"`

	/** 代理 */
	Agent string `json:"agent"          gorm:"column:agent;comment:代理"`

	/** 错误信息 */
	Error string `json:"error"          gorm:"type:longtext;column:error;comment:错误信息"`

	/** 请求Body */
	Body string `json:"body"           gorm:"type:longtext;column:body;comment:请求Body"`

	/** 响应Body */
	Resp string `json:"resp"           gorm:"type:longtext;column:resp;comment:响应Body"`
}
