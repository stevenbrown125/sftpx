// internal/sftp/sftp.go
package sftp

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"sftpx/internal/config"
)

func NewClient(cfg *config.Config) (*sftp.Client, error) {
	auths := []ssh.AuthMethod{}

	// Prefer key auth if provided
	if cfg.SFTP.PrivateKeyPath != "" {
		keyAuth, err := publicKeyAuthFunc(cfg.SFTP.PrivateKeyPath, cfg.SFTP.Passphrase)
		if err != nil {
			return nil, err
		}
		auths = append(auths, keyAuth)
	} else if cfg.SFTP.Password != "" {
		auths = append(auths, ssh.Password(cfg.SFTP.Password))
	}

	sshConfig := &ssh.ClientConfig{
		User:            cfg.SFTP.User,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", cfg.SFTP.Host, cfg.SFTP.Port)
	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, err
	}

	return sftp.NewClient(conn)
}

func publicKeyAuthFunc(keyPath, passphrase string) (ssh.AuthMethod, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	var signer ssh.Signer
	if passphrase == "" {
		signer, err = ssh.ParsePrivateKey(key)
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(passphrase))
	}
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}

// toPOSIX normalizes any OS-built path to POSIX (/), and cleans it.
func toPOSIX(p string) string {
	p = filepath.ToSlash(p)
	p = path.Clean(p)
	// Avoid accidental dir semantics from trailing slash
	if strings.HasSuffix(p, "/") && p != "/" {
		p = strings.TrimRight(p, "/")
	}
	return p
}

// UploadFile uploads a single file and ensures the remote directory exists (Windows -> Linux safe).
func UploadFile(client *sftp.Client, localPath, remotePath string) error {
	src, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("open local file: %w", err)
	}
	defer src.Close()

	// Normalize remote path for POSIX SFTP servers
	rp := toPOSIX(remotePath)

	// Ensure remote directory exists
	remoteDir := path.Dir(rp)
	if remoteDir != "." && remoteDir != "/" {
		if err := client.MkdirAll(remoteDir); err != nil {
			return fmt.Errorf("create remote dir %q: %w", remoteDir, err)
		}
	}

	dst, err := client.Create(rp)
	if err != nil {
		return fmt.Errorf("create remote file %q: %w", rp, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("copy file: %w", err)
	}

	fmt.Printf("Upload complete: %s → %s\n", localPath, rp)
	return nil
}

type Client = sftp.Client
