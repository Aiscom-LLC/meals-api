package swagger

// CateringUserUpdate scheme
type CateringUserUpdate struct {
	FirstName string `json:"firstName,omitempty" example:"Dmitry"`
	LastName  string `json:"lastName,omitempty" example:"Novikov"`
	Email     string `json:"email,omitempty" example:"d.novikov@wellyes.ru"`
}

// UserPasswordUpdate scheme
type UserPasswordUpdate struct {
	OldPassword string `json:"oldPassword" example:"Password12!"`
	NewPassword string `json:"newPassword" example:"Password13!"`
}

// ClientUserUpdate scheme
type ClientUserUpdate struct {
	FirstName string `json:"firstName,omitempty" example:"Dmitry"`
	LastName  string `json:"lastName,omitempty" example:"Novikov"`
	Email     string `json:"email,omitempty" example:"d.novikov@wellyes.ru"`
	Floor     *int   `json:"floor" example:"5"`
	Role      string `json:"role" example:"User"`
}

// RecoveryPassword scheme
type RecoveryPassword struct {
	Email string `json:"email" example:"meals@aisnovations.com"`
}
