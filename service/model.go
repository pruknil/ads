package service

import "time"

const ApplicationId = "772"
const ApplicationABBR = "CCMS"

type ReqHeader struct {
	FuncNm       string `json:"funcNm"`
	RqUID        string `json:"rqUID"`
	RqDt         string `json:"rqDt"`
	RqAppID      string `json:"rqAppId"`
	UserId       string `json:"userId"`
	UserLangPref string `json:"userLangPref"`
	RsUID        string `json:"rsUID,omitempty"`
}

type ResHeader struct {
	FuncNm     string     `json:"funcNm,omitempty"`
	RqUID      string     `json:"rqUID,omitempty"`
	RsAppID    string     `json:"rsAppId,omitempty"`
	RsUID      string     `json:"rsUID,omitempty"`
	RsDt       time.Time  `json:"rsDt,omitempty"`
	StatusCode string     `json:"statusCode,omitempty"`
	ErrorVect  *ErrorVect `json:"errorVect,omitempty"`
}

type ErrorVect struct {
	Error []Error `json:"error"`
}

type Error struct {
	ErrorAppID    string `json:"errorAppId"`
	ErrorAppAbbrv string `json:"errorAppAbbrv"`
	ErrorCode     string `json:"errorCode"`
	ErrorDesc     string `json:"errorDesc"`
	ErrorSeverity string `json:"errorSeverity"`
}

type ReqMsg struct {
	Header ReqHeader   `json:"KBankHeader"`
	Body   interface{} `json:"Body"`
}

type ResMsg struct {
	Header ResHeader   `json:"KBankHeader"`
	Body   interface{} `json:"Body,omitempty"`
}
