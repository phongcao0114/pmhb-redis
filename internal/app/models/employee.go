package models

type (
	// Employee structure
	Employee struct {
		Name     string `json:"name"`
		Position string `json:"position"`
	}
	// GetEmployeeReq structure
	GetEmployeeReq struct {
		Key string `json:"key"`
	}
	// SetEmployeeReq structure
	SetEmployeeReq struct {
		Key        string   `json:"key"`
		Employee   Employee `json:"employee"`
		ExpiryTime int      `json:"expiry_time"`
	}
)
