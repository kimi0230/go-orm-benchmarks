package bench

import (
	"testing"

	"github.com/efectn/go-orm-benchmarks/helper"
	gormV1 "github.com/jinzhu/gorm" // 使用 GORM v1
	_ "github.com/lib/pq"           // 引入 PostgreSQL 驅動
)

type GormV1 struct {
	helper.ORMInterface
	conn *gormV1.DB
}

func CreateGormV1() helper.ORMInterface {
	return &GormV1{}
}

func (gorm *GormV1) Name() string {
	return "gormV1"
}

func (gorm *GormV1) Init() error {
	var err error
	gorm.conn, err = gormV1.Open("postgres", helper.OrmSource) // 使用 GORM v1 的 Open 方法
	if err != nil {
		return err
	}

	gorm.conn.LogMode(false)            // 關閉日誌模式，相當於 v2 的 Logger.Silent
	gorm.conn.DB().SetMaxOpenConns(100) // 設定資料庫連線最大開啟數量
	gorm.conn.DB().SetMaxIdleConns(10)  // 設定資料庫連線最大空閒數量
	gorm.conn.SingularTable(true)       // 如果需要禁用複數表名，這個選項是 GORM v1 的特色
	return nil
}

func (gorm *GormV1) Close() error {
	return gorm.conn.Close() // GORM v1 使用 Close 方法直接關閉連線
}

func (gorm *GormV1) Insert(b *testing.B) {
	m := NewModel()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Id = 0
		err := gorm.conn.Create(m).Error
		if err != nil {
			helper.SetError(b, gorm.Name(), "Insert", err.Error())
		}
	}
}

// Gormv1 不支援多次插入
func (gorm *GormV1) InsertMulti(b *testing.B) {
	ms := make([]*Model, 0, 100)
	for i := 0; i < 100; i++ {
		ms = append(ms, NewModel())
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, m := range ms {
			m.Id = 0
			err := gorm.conn.Create(&m).Error
			if err != nil {
				helper.SetError(b, gorm.Name(), "InsertMulti", err.Error())
			}
		}
	}
}

func (gorm *GormV1) Update(b *testing.B) {
	m := NewModel()

	err := gorm.conn.Create(m).Error
	if err != nil {
		helper.SetError(b, gorm.Name(), "Update", err.Error())
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := gorm.conn.Model(m).Updates(m).Error
		if err != nil {
			helper.SetError(b, gorm.Name(), "Update", err.Error())
		}
	}
}

func (gorm *GormV1) Read(b *testing.B) {
	m := NewModel()

	err := gorm.conn.Create(m).Error
	if err != nil {
		helper.SetError(b, gorm.Name(), "Read", err.Error())
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := gorm.conn.Take(m).Error
		if err != nil {
			helper.SetError(b, gorm.Name(), "Read", err.Error())
		}
	}
}

func (gorm *GormV1) ReadSlice(b *testing.B) {
	m := NewModel()
	for i := 0; i < 100; i++ {
		m.Id = 0
		err := gorm.conn.Create(m).Error
		if err != nil {
			helper.SetError(b, gorm.Name(), "ReadSlice", err.Error())
		}
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var models []*Model
		err := gorm.conn.Where("id > ?", 0).Limit(100).Find(&models).Error
		if err != nil {
			helper.SetError(b, gorm.Name(), "ReadSlice", err.Error())
		}
	}
}
