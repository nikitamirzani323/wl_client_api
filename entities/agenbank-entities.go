package entities

type Model_agenbank struct {
	Agenbank_id         int    `json:"agenbank_id"`
	Agenbank_tipe       string `json:"agenbank_tipe"`
	Agenbank_idbanktype string `json:"agenbank_idbanktype"`
	Agenbank_norek      string `json:"agenbank_norek"`
	Agenbank_nmrek      string `json:"agenbank_nmrek"`
	Agenbank_status     string `json:"agenbank_status"`
	Agenbank_status_css string `json:"agenbank_status_css"`
	Agenbank_create     string `json:"agenbank_create"`
	Agenbank_update     string `json:"agenbank_update"`
}
type Model_agenbankshare struct {
	Agenbank_id   int    `json:"agenbank_id"`
	Agenbank_info string `json:"agenbank_info"`
}

type Controller_agenbanksave struct {
	Page                string `json:"page" validate:"required"`
	Sdata               string `json:"sdata" validate:"required"`
	Agenbank_id         int    `json:"agenbank_id"`
	Agenbank_tipe       string `json:"agenbank_tipe" validate:"required"`
	Agenbank_idbanktype string `json:"agenbank_idbanktype" validate:"required"`
	Agenbank_norek      string `json:"agenbank_norek" validate:"required"`
	Agenbank_nmrek      string `json:"agenbank_nmrek" validate:"required"`
	Agenbank_status     string `json:"agenbank_status" validate:"required"`
}
