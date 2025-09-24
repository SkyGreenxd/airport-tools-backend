package domain

type ToolSetItem struct {
	ToolSetId  int64
	ToolTypeId int64
}

func NewToolSetItem(toolSetId, toolTypeId int64) *ToolSetItem {
	return &ToolSetItem{
		ToolSetId:  toolSetId,
		ToolTypeId: toolTypeId,
	}
}
