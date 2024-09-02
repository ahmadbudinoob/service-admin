package requests

type ResetPassword struct {
	LoginID  string `json:"login_id"`
	Password string `json:"password"`
}

type ResetPin struct {
	LoginID string `json:"login_id"`
	Pin     string `json:"pin"`
}

type LoginRequest struct {
	LoginID  string `json:"loginID"`
	Password string `json:"password"`
}
