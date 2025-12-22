package domain

type EventApplicationStatusChanged struct {
	ApplicationID int64             `json:"application_id"`
	Status        ApplicationStatus `json:"status"`
}
