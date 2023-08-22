package entities

type Model_dpwd struct {
	Dpwd_id           string  `json:"dpwd_id"`
	Dpwd_idmasteragen string  `json:"dpwd_idmasteragen"`
	Dpwd_idmaster     string  `json:"dpwd_idmaster"`
	Dpwd_yearmonth    string  `json:"dpwd_yearmonth"`
	Dpwd_datedoc      string  `json:"dpwd_datedoc"`
	Dpwd_idcurr       string  `json:"dpwd_idcurr"`
	Dpwd_tipedoc      string  `json:"dpwd_tipedoc"`
	Dpwd_tipeakun     string  `json:"dpwd_tipeakun"`
	Dpwd_idagenmember string  `json:"dpwd_idagenmember"`
	Dpwd_multiplier   float32 `json:"dpwd_multiplier"`
	Dpwd_amount       float32 `json:"dpwd_amount"`
	Dpwd_before       float32 `json:"dpwd_before"`
	Dpwd_after        float32 `json:"dpwd_after"`
	Dpwd_ipaddress    string  `json:"dpwd_ipaddress"`
	Dpwd_note         string  `json:"dpwd_note"`
	Dpwd_status       string  `json:"dpwd_status"`
	Dpwd_status_css   string  `json:"dpwd_status_css"`
	Dpwd_create       string  `json:"dpwd_create"`
	Dpwd_update       string  `json:"dpwd_update"`
}

type Controller_dpwdsave struct {
	Page          string `json:"page" validate:"required"`
	Sdata         string `json:"sdata" validate:"required"`
	Dpwd_id       string `json:"dpwd_id"`
	Dpwd_username string `json:"dpwd_username" validate:"required"`
	Dpwd_password string `json:"dpwd_password"`
	Dpwd_name     string `json:"dpwd_name" validate:"required"`
	Dpwd_phone    string `json:"dpwd_phone"`
	Dpwd_email    string `json:"dpwd_email"`
	Dpwd_status   string `json:"dpwd_status"`
}
