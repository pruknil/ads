package http

type KBankRequestHeader struct {
	FuncNm  string `json:"funcNm"`
	RqUID   string `json:"rqUID"`
	RqDt    string `json:"rqDt"`
	RqAppId string `json:"rqAppId"`
	UserId  string `json:"userId"`
}

type KBankResponseHeader struct {
	FuncNm     string `json:"funcNm"`
	RqUID      string `json:"rqUID"`
	RsDt       string `json:"rqDt"`
	RsAppId    string `json:"rqAppId"`
	StatusCode string `json:"statusCode"`
}

type DecryptDataRequest struct {
	KBankRequestHeader     `json:"KBankHeader"`
	DecryptDataBodyRequest `json:"Body"`
}

type DecryptDataResponse struct {
	KBankResponseHeader     `json:"KBankHeader"`
	DecryptDataBodyResponse `json:"Body"`
}

type EncryptDataRequest struct {
	KBankRequestHeader     `json:"KBankHeader"`
	EncryptDataBodyRequest `json:"Body"`
}

type EncryptDataResponse struct {
	KBankResponseHeader     `json:"KBankHeader"`
	EncryptDataBodyResponse `json:"Body"`
}

type EncryptDataBodyRequest struct {
	DPKName string `json:"dPKName"`
	Data    string `json:"data"`
}

type EncryptDataBodyResponse struct {
	DPKName string `json:"dPKName"`
	EData   string `json:"eData"`
}

type DecryptDataBodyRequest struct {
	DPKName string `json:"dPKName"`
	EData   string `json:"eData"`
}

type DecryptDataBodyResponse struct {
	DPKName string `json:"dPKName"`
	Data    string `json:"data"`
}

type SFTP0002I01Request struct {
	KBankRequestHeader     `json:"KBankHeader"`
	SFTP0002I01BodyRequest `json:"Body"`
}

type SFTP0002I01Response struct {
	KBankResponseHeader     `json:"KBankHeader"`
	SFTP0002I01BodyResponse `json:"Body"`
}

type SFTP0002I01BodyRequest struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type SFTP0002I01BodyResponse struct {
}
