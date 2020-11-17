package service

import "github.com/pruknil/ads/backends/http"

type IHttpBackend interface {
	DecryptData(http.DecryptDataRequest) (*http.DecryptDataResponse, error)
	EncryptData(http.EncryptDataRequest) (*http.EncryptDataResponse, error)
	SFTP0002I01(http.SFTP0002I01Request) (*http.SFTP0002I01Response, error)
}
