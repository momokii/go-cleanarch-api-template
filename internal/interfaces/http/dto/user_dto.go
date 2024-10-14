package dto

type UserCreate struct {
	Username string `json:"username" validate:"required,min=5,max=50,alphanum"`
	Password string `json:"password" validate:"required,min=6,max=50,containsany=1234567890,containsany=QWERTYUIOPASDFGHJKLZXCVBNM"`
	Role     int    `json:"role"`
}

type UserUpdate struct {
	Id       int    `json:"id" validate:"required"`
	Username string `json:"username" validate:"required,min=5,max=50,alphanum"`
	Role     int    `json:"role"`
}

type UserChangePassword struct {
	Id          int    `json:"id" validate:"required"`
	OldPassword string `json:"old_password" validate:"required"`
	Password    string `json:"password" validate:"required,min=6,max=50,containsany=1234567890,containsany=QWERTYUIOPASDFGHJKLZXCVBNM"`
}

type UserResponse struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Role      int    `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserSession struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
}
