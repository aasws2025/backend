package model

type UserAccount struct {
	ID       string `json:"userid"`
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Telfon   string `json:"telfon"`
	Alamat   string `json:"alamat"`
	Password string `json:"password"`
}

type AuthRequest struct {
	Email    string `json:"email" valid:"required"`
	Password string `json:"password" valid:"required"`
}
