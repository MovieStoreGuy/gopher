package types

import (
	"github.com/Netflix/go-env"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

type Profile struct {
	GoPath   string `json:"gopath" yaml:"gopath" env:"GOPATH"`
	Name     string `json:"name" yaml:"name" env:"GOPHER_PROFILE_NAME"`
	VCS      string `json:"vcs" yaml:"vcs" env:"GOPHER_VCS"`
	UserName string `json:"username" yaml:"username" env:"GOPHER_USERNAME"`
}

// ConfigureProfile will create a struct using environment variables first
// Then overriding values within the struct from the profile yaml file stored in the ProfilePath
func ConfigureProfile(ProfilePath string) (*Profile, error) {
	var (
		profile Profile
	)
	if _, err := env.UnmarshalFromEnviron(&profile); err != nil {
		return nil, err
	}
	buff, err := ioutil.ReadFile(ProfilePath)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(buff, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// WriteProfile takes a defined profile and writes it to disk
func WriteProfile(ProfilePath string, Profile Profile) error {
	directory := path.Dir(ProfilePath)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, os.ModeDir | 0751); err != nil {
			return err
		}
	}
	buff, err := yaml.Marshal(&Profile)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ProfilePath, buff, 0644)
}
