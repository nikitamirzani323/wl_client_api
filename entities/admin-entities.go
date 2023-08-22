package entities

type Model_admin struct {
	Admin_id           string `json:"admin_id"`
	Admin_idrule       int    `json:"admin_idrule"`
	Admin_tipe         string `json:"admin_tipe"`
	Admin_nmrule       string `json:"admin_nmrule"`
	Admin_username     string `json:"admin_username"`
	Admin_nama         string `json:"admin_nama"`
	Admin_phone1       string `json:"admin_phone1"`
	Admin_phone2       string `json:"admin_phone2"`
	Admin_lastlogin    string `json:"admin_lastlogin"`
	Admin_lastipaddres string `json:"admin_lastipaddres"`
	Admin_status       string `json:"admin_status"`
	Admin_status_css   string `json:"admin_status_css"`
	Admin_create       string `json:"admin_create"`
	Admin_update       string `json:"admin_update"`
}
type Model_adminrule struct {
	Adminrule_idruleadmin int    `json:"adminrule_idruleadmin"`
	Adminrule_nmruleadmin string `json:"adminrule_nmruleadmin"`
}
type Model_adminsave struct {
	Username string `json:"admin_username"`
	Nama     string `json:"admin_nama"`
	Rule     string `json:"admin_rule"`
	Status   string `json:"admin_status"`
	Create   string `json:"admin_create"`
	Update   string `json:"admin_update"`
}

type Controller_adminsave struct {
	Sdata          string `json:"sdata" validate:"required"`
	Page           string `json:"page" validate:"required"`
	Admin_id       string `json:"admin_id"`
	Admin_idrule   int    `json:"admin_idrule" validate:"required"`
	Admin_username string `json:"admin_username" validate:"required"`
	Admin_password string `json:"admin_password"`
	Admin_nama     string `json:"admin_nama" validate:"required"`
	Admin_phone1   string `json:"admin_phone1" validate:"required"`
	Admin_phone2   string `json:"admin_phone2"`
	Admin_status   string `json:"admin_status"`
}

type Responseredis_adminhome struct {
	Admin_username     string `json:"admin_username"`
	Admin_nama         string `json:"admin_nama"`
	Admin_rule         string `json:"admin_rule"`
	Admin_joindate     string `json:"admin_joindate"`
	Admin_timezone     string `json:"admin_timezone"`
	Admin_lastlogin    string `json:"admin_lastlogin"`
	Admin_lastipaddres string `json:"admin_lastipaddres"`
	Admin_status       string `json:"admin_status"`
}
type Responseredis_adminrule struct {
	Adminrule_idrule string `json:"adminrule_idruleadmin"`
}
