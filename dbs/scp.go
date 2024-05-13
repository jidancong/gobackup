package dbs

import (
	"fmt"
	"gobackup/utils"
	"os"
	"path"
	"path/filepath"

	"github.com/povsister/scp"
)

type Scp struct {
	host     string
	port     string
	username string
	password string
}

func NewScp(host, port, username, password string) *Scp {
	return &Scp{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (s *Scp) Backup(fromPath, toPath string) (string, error) {
	config := scp.NewSSHConfigFromPassword(s.username, s.password)

	addr := s.host
	if s.port != "" {
		addr = fmt.Sprintf("%s:%s", s.host, s.port)
	}
	c, err := scp.NewClient(addr, config, &scp.ClientOption{})
	if err != nil {
		return "", err
	}
	defer c.Close()

	localPath := filepath.Join(toPath, fmt.Sprintf("%s-%s-%s-%s", path.Base(fromPath), s.host, utils.TimeFormat(), utils.UUID()))
	if err := os.MkdirAll(localPath, os.ModePerm); err != nil {
		return "", err
	}

	err = c.CopyDirFromRemote(fromPath, localPath, &scp.DirTransferOption{
		PreserveProp: true,
	})
	return localPath, err
}
