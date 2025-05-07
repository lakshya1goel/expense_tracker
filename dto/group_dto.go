package dto

type CreateChatDto struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Users       []string `json:"users" binding:"required"`
}

type AddUsersDto struct {
	Users []string `json:"users" binding:"required"`
}

type UpdateGroupDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
