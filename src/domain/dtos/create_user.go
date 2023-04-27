package dtos

type CreateUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	ImageURL string `json:"image_url"`
}
