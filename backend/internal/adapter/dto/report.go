package dto

type ReportRequest struct {
	ReportType string                 `json:"report_type"`
	Params     map[string]interface{} `json:"params"`
}

type ReportResponse struct {
	ReportID string `json:"report_id"`
	URL      string `json:"url"`
	Status   string `json:"status"`
}
