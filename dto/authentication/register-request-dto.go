package authentication

type RegisterRequestDTO struct {
	Name     string `json:"name" form:"name" binding:"required"`
	Email    string `json:"email" form:"email" binding:"email,required"`
	Password string `json:"password" form:"password" binding:"required"`
}
