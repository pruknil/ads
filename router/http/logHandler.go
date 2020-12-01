package http

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pruknil/ads/logger"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
)

func LogRequest(log *logger.AppLog) gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
		log.Router.WithField(logger.RSUID, c.Writer.Header().Get("X-Request-Id")).Error("[REQ]\n", readBody(rdr1))
		c.Request.Body = rdr2
		c.Next()
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogResponse(log *logger.AppLog) gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		if !strings.HasSuffix(c.Request.RequestURI, "jpg") {
			log.Router.WithField(logger.RSUID, c.Writer.Header().Get("X-Request-Id")).Error("[RES]\n", maskingValue(blw.body.String()))
		}
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(reader)
	re := regexp.MustCompile(`[\t\n\r\s]`)
	s := re.ReplaceAllString(buf.String(), "")
	return maskingValue(s)
}

func maskingValue(input string) string {
	maskField := []string{"cardNumber", "cvv"}
	for _, s := range maskField {
		re, _ := regexp.Compile(fmt.Sprintf(`"%s":"(\w+)"`, s))
		resultSlice := re.FindStringSubmatch(input)
		if len(resultSlice) > 0 {
			input = strings.Replace(input, resultSlice[0], fmt.Sprintf(`"%s":"***"`, s), -1)
		}
	}
	return input
}
