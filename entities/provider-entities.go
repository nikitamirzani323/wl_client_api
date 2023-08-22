package entities

type Model_provider struct {
	Provider_id         int    `json:"provider_id"`
	Provider_name       string `json:"provider_name"`
	Provider_owner      string `json:"provider_owner"`
	Provider_email      string `json:"provider_email"`
	Provider_phone1     string `json:"provider_phone1"`
	Provider_phone2     string `json:"provider_phone2"`
	Provider_urlwebsite string `json:"provider_urlwebsite"`
	Provider_status     string `json:"provider_status"`
	Provider_status_css string `json:"provider_status_css"`
	Provider_create     string `json:"provider_create"`
	Provider_update     string `json:"provider_update"`
}
type Model_providershare struct {
	Provider_id   int    `json:"provider_id"`
	Provider_name string `json:"provider_name"`
}
type Controller_providersave struct {
	Page                string `json:"page" validate:"required"`
	Sdata               string `json:"sdata" validate:"required"`
	Provider_id         int    `json:"provider_id"`
	Provider_name       string `json:"provider_name" validate:"required"`
	Provider_owner      string `json:"provider_owner" validate:"required"`
	Provider_email      string `json:"provider_email"`
	Provider_phone1     string `json:"provider_phone1" validate:"required"`
	Provider_phone2     string `json:"provider_phone2"`
	Provider_urlwebsite string `json:"provider_urlwebsite"`
	Provider_status     string `json:"provider_status" validate:"required"`
}
