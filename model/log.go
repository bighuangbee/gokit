package model


type SysLog struct {
	SysLogRpc
	// 不需要赋值
	CreateSecond int64 `json:"-" gorm:"create_second;type:bigint not null;comment:创建时间戳"`
	// 不需要赋值
	CreateAt MyTime `json:"createAt"  gorm:"create_at"`
}
type SysLogRpc struct {
	Id        int64        `json:"id"` // 自增id
	CorpId    uint64       `json:"corpId" gorm:"corp_id"`
	UserId    uint32       `json:"userId" gorm:"user_id"`
	Name      string       `json:"name" gorm:"name"`
	LoginName string       `json:"account" gorm:"login_name"`
	Type      LogType 		`json:"type" gorm:"type"`
	IP        string       `json:"ip" gorm:"ip"`
	Url       string       `json:"url" gorm:"url"`
	UrlTitle  string       `json:"urlTitle" gorm:"url_title"`
	Params    string       `json:"params" gorm:"params"`
	UserAgent string       `json:"userAgent" gorm:"_"`
}


type LogType uint8

const (
	LogDefault LogType = iota
	// 创建
	LogCreate
	// 查看
	LogView
	// 编辑
	LogEdit
	// 删除
	LogDelete
	// 登陆
	LogLogin
	// 登出
	LogLogout
	// 导出
	LogExport
	// 导入
	LogImport
	// 保存
	LogSave
)
