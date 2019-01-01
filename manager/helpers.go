package manager

import (
	"github.com/MovieStoreGuy/gopher/types"
	"os"
	"path"
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
