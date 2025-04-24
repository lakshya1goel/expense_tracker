package dto

type CreateGroupDto struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Members     []string `json:"members" binding:"required"`
}