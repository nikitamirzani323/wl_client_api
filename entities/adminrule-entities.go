package entities

type Model_agenadminrule struct {
	Agenadminrule_id     int    `json:"agenadminrule_id"`
	Agenadminrule_name   string `json:"agenadminrule_name"`
	Agenadminrule_rule   string `json:"agenadminrule_rule"`
	Agenadminrule_create string `json:"agenadminrule_create"`
	Agenadminrule_update string `json:"agenadminrule_update"`
}
type Controller_agenadminrulesave struct {
	Sdata              string `json:"sdata" validate:"required"`
	Page               string `json:"page" validate:"required"`
	Agenadminrule_id   int    `json:"agenadminrule_id" `
	Agenadminrule_name string `json:"agenadminrule_name"`
	Agenadminrule_rule string `json:"agenadminrule_rule"`
}
