package dao

type PersonDaoInterface interface {
	BaseDaoInterface
}

type PersonDao struct {
	BaseDao `inject:"inline"`
}
