package model

type CreateRequestRecordModel struct {
	RequestBody string `json:"request_body"`
}

type GetResponseModel struct {
	RequestId string `json:"request_id"`
}

type CalculationResponseModel struct {
	RequestId   string `json:"request_id"`
	RequestBody string `json:"request_body"`
	Response    string `json:"response"`
	Created_at  string `json:"created_at"`
	Updated_at  string `json:"updated_at"`
}
