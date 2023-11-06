package form

type RegisterForm struct {
	Name     string `json:"name,omitempty"`
	Mobile   string `json:"mobile,omitempty"`
	Password string `json:"password,omitempty"`
	SMSCode  string `json:"sms-code,omitempty"`
}

type LoginForm struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}
