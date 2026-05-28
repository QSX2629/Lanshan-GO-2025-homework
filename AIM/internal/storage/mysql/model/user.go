package model

type User struct {
	ID       string `gorm:"column:id;primaryKey"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	Avatar   string `gorm:"column:avatar"`
}

func (u User) TableName() string {
	return "user"
}
