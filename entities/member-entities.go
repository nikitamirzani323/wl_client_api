package entities

type Model_member struct {
	Member_id         string      `json:"member_id"`
	Member_username   string      `json:"member_username"`
	Member_timezone   string      `json:"member_timezone"`
	Member_ipaddress  string      `json:"member_ipaddress"`
	Member_lastlogin  string      `json:"member_lastlogin"`
	Member_name       string      `json:"member_name"`
	Member_phone      string      `json:"member_phone"`
	Member_email      string      `json:"member_email"`
	Member_listbank   interface{} `json:"member_listbank"`
	Member_status     string      `json:"member_status"`
	Member_status_css string      `json:"member_status_css"`
	Member_create     string      `json:"member_create"`
	Member_update     string      `json:"member_update"`
}
type Model_membershare struct {
	Member_id       string      `json:"member_id"`
	Member_name     string      `json:"member_name"`
	Member_listbank interface{} `json:"member_listbank"`
}
type Model_memberbank struct {
	Memberbank_id          int    `json:"memberbank_id"`
	Memberbank_idbanktype  string `json:"memberbank_idbanktype"`
	Memberbank_norek       string `json:"memberbank_norek"`
	Memberbank_nmownerbank string `json:"memberbank_nmownerbank"`
}
type Model_memberbankshare struct {
	Memberbank_id   int    `json:"memberbank_id"`
	Memberbank_info string `json:"memberbank_info"`
}

type Controller_membersave struct {
	Page            string `json:"page" validate:"required"`
	Sdata           string `json:"sdata" validate:"required"`
	Member_id       string `json:"member_id"`
	Member_username string `json:"member_username" validate:"required"`
	Member_password string `json:"member_password"`
	Member_name     string `json:"member_name" validate:"required"`
	Member_phone    string `json:"member_phone"`
	Member_email    string `json:"member_email"`
	Member_status   string `json:"member_status"`
}
type Controller_memberbanksave struct {
	Page                    string `json:"page" validate:"required"`
	Sdata                   string `json:"sdata" validate:"required"`
	Memberbank_idagenmember string `json:"memberbank_idagenmember" validate:"required"`
	Memberbank_idbanktype   string `json:"memberbank_idbanktype" validate:"required"`
	Memberbank_norek        string `json:"memberbank_norek" validate:"required"`
	Memberbank_nmownerbank  string `json:"memberbank_nmownerbank" validate:"required"`
}
type Controller_memberbankdelete struct {
	Page                    string `json:"page" validate:"required"`
	Sdata                   string `json:"sdata" validate:"required"`
	Memberbank_id           int    `json:"memberbank_id" validate:"required"`
	Memberbank_idagenmember string `json:"memberbank_idagenmember" validate:"required"`
}
type Controller_membersharesearch struct {
	Page   string `json:"page" validate:"required"`
	Sdata  string `json:"sdata" validate:"required"`
	Search string `json:"search" `
}
