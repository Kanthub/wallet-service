package backend

type AdminLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AdminLoginResponse struct {
	Success   bool       `json:"success"`
	Message   string     `json:"message"`
	Token     string     `json:"token,omitempty"`
	AdminInfo *AdminInfo `json:"user_info,omitempty"`
}

type AdminInfo struct {
	Guid     string `json:"guid"`
	Username string `json:"username"`
}

type AdminLogoutRequest struct {
	UserId string `json:"user_id" binding:"required"`
	Token  string `json:"token"`
}

type AdminLogoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
