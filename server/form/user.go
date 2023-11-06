package form

type RegisterForm struct {
	Name     string `json:"name,omitempty"  binding:"required"`
	Mobile   string `json:"mobile,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required,min=3,max=20"`
	SMSCode  string `json:"sms-code,omitempty" binding:"required,min=6,max=6"`
}

type LoginForm struct {
	Name     string `json:"name,omitempty"  binding:"required"`
	Password string `json:"password,omitempty" binding:"required,min=3,max=20"`
}
