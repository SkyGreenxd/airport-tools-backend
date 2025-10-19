package domain

const (
	Engineer       string = "Engineer"
	QualityAuditor string = "Quality Auditor"
)

type Role struct {
	Id   int64
	Name string

	Users []*User
}

func NewRole(id int64, name string) *Role {
	return &Role{
		Id:   id,
		Name: name,
	}
}
