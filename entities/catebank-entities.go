package entities

type Model_catebank struct {
	Catebank_id         int         `json:"catebank_id"`
	Catebank_name       string      `json:"catebank_name"`
	Catebank_status     string      `json:"catebank_status"`
	Catebank_status_css string      `json:"catebank_status_css"`
	Catebank_list       interface{} `json:"catebank_list"`
	Catebank_create     string      `json:"catebank_create"`
	Catebank_update     string      `json:"catebank_update"`
}
type Model_bankType struct {
	Banktype_id         string `json:"banktype_id"`
	Banktype_name       string `json:"banktype_name"`
	Banktype_img        string `json:"banktype_img"`
	Banktype_status     string `json:"banktype_status"`
	Banktype_status_css string `json:"banktype_status_css"`
	Banktype_create     string `json:"banktype_create"`
	Banktype_update     string `json:"banktype_update"`
}
type Model_bankTypeshare struct {
	Catebank_name string `json:"catebank_name"`
	Banktype_id   string `json:"banktype_id"`
}

type Controller_catebanksave struct {
	Page            string `json:"page" validate:"required"`
	Sdata           string `json:"sdata" validate:"required"`
	Catebank_id     int    `json:"catebank_id"`
	Catebank_name   string `json:"catebank_name" validate:"required"`
	Catebank_status string `json:"catebank_status" validate:"required"`
}
type Controller_banktypesave struct {
	Page                string `json:"page" validate:"required"`
	Sdata               string `json:"sdata" validate:"required"`
	Banktype_id         string `json:"banktype_id" validate:"required"`
	Banktype_idcatebank int    `json:"banktype_idcatebank" validate:"required"`
	Banktype_name       string `json:"banktype_name" validate:"required"`
	Banktype_img        string `json:"banktype_img" validate:"required"`
	Banktype_status     string `json:"banktype_status" validate:"required"`
}
