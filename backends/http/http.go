package http

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	breaker "github.com/sony/gobreaker"
	"io/ioutil"
	"github.com/pruknil/ads/logger"
	"net/http"
	"time"
)

type Client struct {
	*breaker.CircuitBreaker
	Config
	logger.AppLog
	RestRequestBuilder
}

func New(c Config, log logger.AppLog) IHttpBackendService {
	var st breaker.Settings
	st.Name = "HTTP"
	st.Timeout = 3
	st.ReadyToTrip = func(counts breaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}
	cb := breaker.NewCircuitBreaker(st)
	return &Client{CircuitBreaker: cb, Config: c, AppLog: log}
}

type CCMSAPI struct {
	DecryptDataUrl string
	EncryptDataUrl string
	LocalUrl       string
}
type Config struct {
	CCMSAPI CCMSAPI
}

func (s *Client) DoRequest(input RestRequestBuilder) error {
	defer func(t time.Time) {
		s.AppLog.Trace.Printf("%s", input.Req.URL)
		s.AppLog.Trace.Printf("http.DoRequest elapsed time %.4f ms", float64(time.Since(t).Nanoseconds())/float64(time.Millisecond))

		//s.Perf.WithFields(logrus.Fields{
		//	logger.FUNCNM: "",
		//	logger.RQUID:  "",
		//	logger.RSUID:  "",
		//	logger.STATUS: "",
		//	logger.TM:     fmt.Sprintf("%.4f", float64(time.Since(t).Nanoseconds())/float64(time.Millisecond)),
		//}).Print()

	}(time.Now())
	res, err := s.CircuitBreaker.Execute(func() (i interface{}, e error) {

		rawBody, err := input.Req.GetBody()
		if err != nil {
			return nil, err
		}
		body, err := ioutil.ReadAll(rawBody)
		if err != nil {
			return nil, err
		}
		s.Rest.Println("REQ Body:", string(body))
		resp, err := input.client.Do(input.Req) //client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		s.Rest.Println("RES Body:", string(body))
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return body, fmt.Errorf("%d:%s", resp.StatusCode, http.StatusText(resp.StatusCode))
		}
		return body, nil
	})
	if res != nil {
		errParse := json.Unmarshal(res.([]byte), input.responseModel)
		if errParse != nil {
			return errParse
		}
		if err != nil {
			return err
		}
	}
	return err
}

/*type Req struct {
	*http.Request
	Url    string
	Method string
	Body   io.Reader
	Header map[string][]string
}*/

type RestRequestBuilder struct {
	Req           *http.Request
	client        *http.Client
	responseModel interface{}
}

func (builder *RestRequestBuilder) SetQueryParams(query map[string]string) *RestRequestBuilder {
	if query != nil {
		q := builder.Req.URL.Query()
		for k, v := range query {
			q.Add(k, v)
		}
		builder.Req.URL.RawQuery = q.Encode()
	}
	return builder
}

func (builder *RestRequestBuilder) SetBasicAuth(username string, pass string) *RestRequestBuilder {
	builder.Req.SetBasicAuth(username, pass)
	return builder
}

func (builder *RestRequestBuilder) SetNormalTransport() *RestRequestBuilder {
	builder.client = &http.Client{}
	return builder
}

func (builder *RestRequestBuilder) SetSecurityTransport() *RestRequestBuilder {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	builder.client = &http.Client{
		Transport: transport,
	}
	return builder
}

func (builder *RestRequestBuilder) SetAuth(token string) *RestRequestBuilder {
	if token != "" {
		builder.Req.Header.Set("Authorization", token)
	}
	return builder
}

func (builder *RestRequestBuilder) SetHeader(header map[string]string) *RestRequestBuilder {
	if header != nil {
		for k, v := range header {
			builder.Req.Header.Set(k, v)
		}

	}
	return builder
}

func (builder *RestRequestBuilder) SetContentType(contentType string) *RestRequestBuilder {
	if contentType == "" {
		contentType = "application/json"
	}
	builder.Req.Header.Set("Content-type", contentType)
	return builder
}

func (builder *RestRequestBuilder) SetResponseModel(model interface{}) *RestRequestBuilder {
	builder.responseModel = model
	return builder
}
