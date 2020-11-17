package sftp

import (
	"fmt"
	"github.com/taruti/sshutil"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"github.com/pruknil/ads/logger"
	"github.com/pruknil/ads/router/sftpd"
	"log"
	"time"
)

type Config struct {
	Host string
	Port string

	UserId   string
	Password string

	ConnTimeout   time.Duration
	ReadDeadline  time.Duration
	WriteDeadline time.Duration
	PrivateKey    string
}

func NewSftp(cfg Config, log logger.AppLog, ccmsFs *CcmsFs) *SFTP {
	return &SFTP{
		config: cfg,
		log:    log,
		fs:     ccmsFs,
	}
}

type SFTP struct {
	config Config
	log    logger.AppLog
	fs     *CcmsFs
}

func (s *SFTP) Start() {
	fmt.Println("Start sftp")
	go s.RunServer(fmt.Sprintf("%s:%s", s.config.Host, s.config.Port), s.fs)
}

func (s *SFTP) Shutdown() {
	fmt.Println("Shutdown sftp")
}

func (s *SFTP) RunServer(hostport string, fs sftpd.FileSystem) {
	var testUser = s.config.UserId
	var testPass = []byte(s.config.Password)
	cfg := sftpd.Config{HostPort: hostport, FileSystem: fs, LogFunc: log.Println}
	cfg.Init()
	cfg.PasswordCallback = sshutil.CreatePasswordCheck(testUser, testPass)

	// Add the sshutil.RSA2048 and sshutil.Save flags if needed for the server in question...
	hkey, e := sshutil.KeyLoader{Flags: sshutil.Create}.Load()
	if e != nil {
		log.Println(e)
		return
	}

	privateBytes, err := ioutil.ReadFile(s.config.PrivateKey)
	if err != nil {
		log.Fatal("Failed to load private key", err)
	}

	hkey, err = ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key", err)
	}

	cfg.AddHostKey(hkey)

	log.Printf("Listening on %s user %s pass %s\n", hostport, testUser, testPass)
	cfg.RunServer()
}
