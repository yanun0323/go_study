package entity

type Model struct {
	ID    int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Value string `gorm:"column:value"`
}

func (Model) TableName() string {
	return "models"
}
