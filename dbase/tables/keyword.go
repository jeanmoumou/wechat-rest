package tables

// 关键词

type Keyword struct {
	Rd        uint   `json:"rd" gorm:"primaryKey"`                           // 主键
	Group     string `json:"group" gorm:"uniqueIndex:t_r_phrase"`            // 分组
	Roomid    string `json:"roomid" gorm:"uniqueIndex:t_r_phrase;default:-"` // 群聊 id
	Phrase    string `json:"phrase" gorm:"uniqueIndex:t_r_phrase"`           // 短语
	Target    string `json:"target"`                                         // 目标
	Level     int32  `json:"level" gorm:"default:-1"`                        // 等级
	CreatedAt int64  `json:"created_at"`                                     // 创建时间戳
	UpdatedAt int64  `json:"updated_at"`                                     // 最后更新时间戳
}
