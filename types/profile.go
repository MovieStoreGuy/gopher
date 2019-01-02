package types

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"

	"github.com/Netflix/go-env"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

// Profile defines how a user would usual develop under their given profile
// This enables saving time by fetch / storing data instead of asking for it each call
type Profile struct {
	GoPath        string `json:"gopath" yaml:"gopath" env:"GOPATH"`
	Name          string `json:"name" yaml:"name" env:"GOPHER_PROFILE_NAME"`
	VCS           string `json:"vcs" yaml:"vcs" env:"GOPHER_VCS" validate:"required"`
	UserName      string `json:"username" yaml:"username" env:"GOPHER_USERNAME" validate:"required"`
	EnableModules string `json:"enableModules" yaml:"enableModules" env:"GO111MODULE"`

	SshKeyPath string `json:"sshkeypath" yaml:"sshkeypath" env:"GOPHER_SSH_KEYPATH"`
	SshKey     string `json:"-" yaml:"-" env:"GOPHER_SSH_KEY"`
}

func (p *Profile) loadSshkey() error {
	if p.SshKey != "" || p.SshKeyPath == "" {
		return nil
	}
	buff, err := ioutil.ReadFile(p.SshKeyPath)
	if err != nil {
		return err
	}
	p.SshKey = bytes.NewBuffer(buff).String()
	return nil
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ConfigureProfile will create a struct using environment variables first
// Then overriding values within the struct from the profile yaml file stored in the ProfilePath
func ConfigureProfile(ProfilePath string) (*Profile, error) {
	var (
		profile Profile
	)
	buff, err := ioutil.ReadFile(ProfilePath)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(buff, &profile); err != nil {
		return nil, err
	}
	if _, err := env.UnmarshalFromEnviron(&profile); err != nil {
		return nil, err
	}
	return &profile, (&profile).loadSshkey()
}

// WriteProfile takes a defined profile and writes it to disk
func WriteProfile(ProfilePath string, Profile Profile) error {
	directory := path.Dir(ProfilePath)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, os.ModeDir|0751); err != nil {
			return err
		}
	}
	buff, err := yaml.Marshal(&Profile)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ProfilePath, buff, 0644)
}

// ValidateProfile ensures all required fields of the profile are met
// so that other functions may use them to ensure it is correct
func ValidateProfile(Profile *Profile) error {
	return validate.Struct(Profile)
}
