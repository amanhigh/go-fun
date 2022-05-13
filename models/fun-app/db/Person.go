package db

type Person struct {
	Name string `gorm:"not null" binding:"required"`
	Age  int    `gorm:"not null" binding:"required"`

	Gender string `gorm:"not null" binding:"required,eq=MALE|eq=FEMALE"`

	Version int64 `gorm:"not null" json:"-" binding:"-"`
}
