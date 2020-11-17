package sftp

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"github.com/pruknil/ads/backends/http"
	"github.com/pruknil/ads/backends/http/service"
	"github.com/pruknil/ads/router/sftpd"
	"log"
	"time"
)

func NewCcmsFs(httpBackend service.IHttpBackend) *CcmsFs {
	return &CcmsFs{
		IHttpBackend: httpBackend,
		CcmsDir:      &CcmsDir{map[string]*CcmsFile{}},
	}
}

type CcmsFs struct {
	sftpd.EmptyFS
	service.IHttpBackend
	*CcmsDir
}

type CcmsDir struct {
	M map[string]*CcmsFile
}
type FState int

const (
	FILE_PUT FState = iota
	FILE_GET
)

type CcmsFile struct {
	service.IHttpBackend
	sftpd.EmptyFile
	Bs         []byte
	CreateTime time.Time
	User       string
	IsDir      bool
	State      FState
	Name       string
}

func (fs *CcmsFs) Mkdir(path string, attr *sftpd.Attr) error {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	fs.M[path] = &CcmsFile{
		CreateTime: time.Now(),
		User:       "testuser",
		IsDir:      true,
		Name:       path,
	}
	return nil
}

func (fs *CcmsFs) Remove(path string) error {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	if _, ok := fs.M[path]; ok {
		delete(fs.M, path)
	}

	return nil
}

func (fs *CcmsFs) OpenDir(path string) (sftpd.Dir, error) {
	d := CcmsDir{M: map[string]*CcmsFile{}}
	for k, v := range fs.CcmsDir.M {
		d.M[k] = v
	}
	return &d, nil
}

func (fs *CcmsFs) OpenFile(path string, mode uint32, attr *sftpd.Attr) (sftpd.File, error) {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	switch mode {
	case 1:
		f, ok := fs.CcmsDir.M[path]
		if ok {
			return f, nil
		}
		return nil, errors.New("Not found!")
	case 26:
		syf := &CcmsFile{
			IHttpBackend: fs.IHttpBackend,
			CreateTime:   time.Now(),
			User:         "testuser",
			State:        FILE_PUT,
			Name:         path,
		}
		fs.CcmsDir.M[path] = syf
		return syf, nil
	default:
		return nil, errors.New(fmt.Sprintf("Invalid mode %d", mode))
	}
}

func (fs *CcmsFs) Stat(path string, islstat bool) (*sftpd.Attr, error) {
	var a sftpd.Attr
	log.Println("STAT", path)
	if path == "" || path == "/" || path == "." {
		a.Flags = sftpd.ATTR_MODE
		a.Mode = sftpd.MODE_DIR | 0777
		return &a, nil
	}
	if path[0] == '/' {
		path = path[1:]
	}
	f, ok := fs.CcmsDir.M[path]
	if ok {
		f.fillAttr(&a)
		return &a, nil
	}
	return nil, errors.New("not found")
}

func (d *CcmsDir) Readdir(count int) ([]sftpd.NamedAttr, error) {
	var rs []sftpd.NamedAttr
	for k, v := range d.M {
		na := sftpd.NamedAttr{Name: k}
		v.fillAttr(&na.Attr)
		rs = append(rs, na)
		delete(d.M, k)
	}
	if len(rs) == 0 {
		return nil, io.EOF
	}
	return rs, nil
}

func (d *CcmsDir) Close() error {
	return nil
}

func (f *CcmsFile) ReadAt(bs []byte, offset int64) (int, error) {
	f.State = FILE_GET
	return bytes.NewReader(f.Bs).ReadAt(bs, offset)
}

func (f *CcmsFile) WriteAt(bs []byte, offset int64) (int, error) {
	rn := []rune(string(bs))
	bb := bytes.NewBuffer(f.Bs)
	bb.WriteString(string(rn))
	f.Bs = bb.Bytes()
	return len(bs), nil
}

func (f *CcmsFile) fillAttr(attr *sftpd.Attr) {
	attr.Flags = sftpd.ATTR_SIZE | sftpd.ATTR_MODE
	attr.Size = uint64(len(f.Bs))
	if !f.IsDir {
		attr.Mode = sftpd.MODE_REGULAR | 0777
	} else {
		attr.Mode = sftpd.MODE_DIR | 0777
	}

	attr.User = f.User
	attr.MTime = f.CreateTime
}
func (f *CcmsFile) Close() error {
	switch f.State {
	case FILE_GET:
		fmt.Println("Get file Done")
	case FILE_PUT:
		fmt.Println("Put file " + f.Name + " Done")
		f.collectService(f.Name)
	}
	return nil
}

func (f *CcmsFile) collectService(fileName string) {
	switch fileName {
	case "BASE24P.REFUND.FEE.FROMPCB.txt":
		f.IHttpBackend.SFTP0002I01(http.SFTP0002I01Request{
			KBankRequestHeader: http.KBankRequestHeader{
				FuncNm:  "SFTP0002I01",
				RqUID:   uuid.New().String(),
				RqDt:    time.Now().Format(time.RFC3339),
				RqAppId: "772",
				UserId:  "testuser",
			},
			SFTP0002I01BodyRequest: http.SFTP0002I01BodyRequest{
				Input:  f.Name,
				Output: "ENC." + f.Name,
			},
		})
	}

}
