package manager

import (
	"bytes"
	"os"
	"path"

	"github.com/MovieStoreGuy/gopher/types"
	"golang.org/x/crypto/ssh"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

func usingModules(Profile *types.Profile) bool {
	if gopath, exist := os.LookupEnv("GOPATH"); exist && Profile.GoPath == gopath {
		return false
	}
	return Profile.EnableModules == "auto" || Profile.EnableModules == "on"
}

func getProjectRoot(Profile *types.Profile) (string, error) {
	if err := types.ValidateProfile(Profile); err != nil {
		return "", err
	}
	if usingModules(Profile) {
		return Profile.GoPath, nil
	}
	return path.Join(Profile.GoPath, "src", Profile.VCS, Profile.UserName), nil
}

func gitSShAuthConfig(Profile *types.Profile) (gitssh.PublicKeys, error) {
	var (
		config gitssh.PublicKeys
	)
	if err := types.ValidateProfile(Profile); err != nil {
		return config, err
	}
	signer, err := ssh.ParsePrivateKey(bytes.NewBufferString(Profile.SshKey).Bytes())
	if err != nil {
		return config, err
	}
	config = gitssh.PublicKeys{
		User:   Profile.UserName,
		Signer: signer,
	}
	return config, nil
}
