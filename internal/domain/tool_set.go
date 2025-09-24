package domain

type ToolSet struct {
	Id   int64
	Name string

	Tools []*ToolType
}

func NewToolSet(name string) *ToolSet {
	return &ToolSet{
		Name: name,
	}
}
