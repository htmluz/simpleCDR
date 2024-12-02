package models

type GetDaysResponse struct {
	Days      int    `json:"days"`
	UpdatedAt string `json:"updated_at"`
}

type UpdateDaysRequest struct {
	Days int `json:"days"`
}
