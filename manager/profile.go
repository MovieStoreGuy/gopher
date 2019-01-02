package manager

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/MovieStoreGuy/gopher/types"
	"github.com/fatih/color"
)

var (
	// GopherPath defines where to find the saved profiles
	GopherPath = path.Join(os.Getenv("HOME"), ".gopher")
	// CurrentProfile is the profile to be used as a default
	CurrentProfile = path.Join(GopherPath, "current")
)

func init() {
	if _, err := os.Stat(CurrentProfile); os.IsNotExist(err) {
		color.Yellow("Creating a blank profile to use, ensure to create a new one")
		if err := StoreProfile("current", &types.Profile{GoPath: os.Getenv("GOPATH")}); err != nil {
			color.Red("Unable to create a blank profile due to %v", err)
			os.Exit(1)
		}
		color.Yellow("Created a blank profile, be sure to define a profile")
	}
}

// LoadProfile will load the given profile name wrapped by the GopherPath
func LoadProfile(name string) (*types.Profile, error) {
	return types.ConfigureProfile(path.Join(GopherPath, name))
}

// StoreProfile will store the given profile to the name file wrapped by the GopherProfile
func StoreProfile(name string, Profile *types.Profile) error {
	return types.WriteProfile(path.Join(GopherPath, name), *Profile)
}

// SetDefaultProfile will create a symlink using the CurrentProfile to soft link to the
// actual profile to use
func SetDefaultProfile(name string) error {
	if name == "" {
		return errors.New("can not set profile to equla nothing")
	}
	target := path.Join(GopherPath, name)
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return err
	}
	if _, err := os.Stat(CurrentProfile); !os.IsNotExist(err) {
		if err := os.Remove(CurrentProfile); err != nil {
			return err
		}
	}
	return os.Symlink(target, CurrentProfile)
}

// GetStoredProfiles gets all the profiles stored within the GopherPath excluding current
// as it is a soft link to a profile that exists
func GetStoredProfiles() ([]*types.Profile, error) {
	files, err := ioutil.ReadDir(GopherPath)
	if err != nil {
		return nil, err
	}
	var profiles []*types.Profile
	for _, f := range files {
		if !f.IsDir() && (f.Name() != "current" || strings.HasPrefix(f.Name(), ".")) {
			p, err := LoadProfile(f.Name())
			if err != nil {
				return nil, err
			}
			profiles = append(profiles, p)
		}
	}
	return profiles, nil
}
