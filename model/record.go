package model

import (
	"errors"
	"time"
)

// Record 结构体定义了运动活动的数据模型, 只能 insert 或者 delete 数据, 不能更新数据.
type Record struct {
	Model // 包含默认的id, created_at, updated_at等字段
	// User
	UserID string `gorm:"not null;" json:"user_id"`
	// 运动类型字段
	Type string `gorm:"size:20;not null;default:'';comment:'运动类型（跑步、慢走等）'" json:"type"`
	// 起始时间字段
	StartTime time.Time `gorm:"not null;comment:'运动开始时间'" json:"start_time"`
	// 结束时间字段
	EndTime  time.Time `gorm:"not null;comment:'运动结束时间'" json:"end_time"`
	Distance float64   `gorm:"not null;comment:距离" json:"distance"`
	// 平均速度字段（单位可根据实际情况设定，例如千米每小时）
	AverageSpeed float64 `gorm:"not null;default:0;comment:'平均速度 in M/s'" json:"average_speed"`
	// 路线信息字段（可以是地理坐标序列，或者字符串形式的描述）
	RouteData []byte `gorm:"comment:'路线信息'" json:"route_data"` // 使用jsonb类型存储复杂结构数据
	// 如果路线是一系列GPS点，可以定义如下（假设Route为自定义的GPS点数组结构）
	//Route []Route `gorm:"-" json:"route"` // 不直接映射到数据库字段，可通过Marshal/Unmarshal转换为JSON存入RouteData字段
	// Room 信息不保存
}

// Route GPS点结构体
type Route struct {
	Latitude  float64 `gorm:"" json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (r *Record) Insert() error {
	err := DB.Create(r).Error
	return err
}

func (r *Record) Delete() error {
	if r.ID == 0 {
		return errors.New("id 为空！")
	}
	err := DB.Delete(r).Error
	return err
}

func GetMaxRecordID() int {
	var record Record
	DB.Last(&record)
	return int(record.ID)
}

func DeleteRecordByID(id uint) error {
	if id == 0 {
		return errors.New("id 为空！")
	}
	record := Record{Model: Model{ID: id}}
	return record.Delete()
}

func GetRecords(userId uint, num int) ([]*Record, error) {
	var records []*Record
	var err error
	err = DB.Order("id desc").Limit(num).Find(&records).Error
	return records, err
}
