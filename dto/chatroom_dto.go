package dto

type CreateChatDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Members     []string `json:"members" binding:"required"`
}

type AddMembersDto struct {
	Members []string `json:"members" binding:"required"`
}

type UpdateGroupDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}