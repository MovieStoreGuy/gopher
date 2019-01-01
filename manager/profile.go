package manager

import (
	"github.com/MovieStoreGuy/gopher/types"
	"github.com/fatih/color"
	"os"
	"path"
)

var (
	// GopherPath defines where to find the saved profiles
	GopherPath     = path.Join(os.Getenv("HOME"), ".gopher")
	CurrentProfile = path.Join(GopherPath, "current")
)

func init() {
	if _, err := os.Stat(CurrentProfile); os.IsNotExist(err) {
		color.Yellow("Creating a blank profile to use, ensure to create a new one")
		if err := StoreProfile("current", &types.Profile{GoPath: os.Getenv("GOPATH")}); err != nil {
			color.Red("Unable to create a blank profile due to %v", err)
			return
		}
		color.Yellow("Created a blank profile, be sure to define a profile")
	}
}

func LoadProfile(name string) (*types.Profile, error) {
	return types.ConfigureProfile(path.Join(GopherPath, name))
}

func StoreProfile(name string, Profile *types.Profile) error {
	return types.WriteProfile(path.Join(GopherPath, name), *Profile)
}

func SetDefaultProfile(name string) error {
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
