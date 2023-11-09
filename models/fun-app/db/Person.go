package db

type Person struct {
	// Validations - https://gin-gonic.com/docs/examples/binding-and-validation/
	Id   int64  `gorm:"primaryKey"`
	Name string `gorm:"not null" binding:"required,min=1,max=25"`
	Age  int    `gorm:"not null" binding:"required,min=1,max=150"`

	Gender string `gorm:"not null" binding:"required,eq=MALE|eq=FEMALE" enums:"MALE,FEMALE"`

	//TODO: Implement Versioning
	Version int64 `gorm:"not null" json:"-" binding:"-"`
}
