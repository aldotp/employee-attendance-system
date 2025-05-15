package dto

type NotificationRequest struct {
	RecipientID string `json:"recipient_id"`
	Title       string `json:"title"`
	Message     string `json:"message"`
}

type NotificationResponse struct {
	NotificationID string `json:"notification_id"`
	Status         string `json:"status"`
}

type ListNotificationRequest struct {
	Skip  uint64 `json:"skip" form:"skip"`
	Limit uint64 `json:"limit" form:"form"`
}
