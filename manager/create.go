package manager

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/MovieStoreGuy/gopher/types"
	"gopkg.in/src-d/go-git.v4"
)

const (
	DefaultReadMeTemplate = `# %s`
)

func CreateProject(Profile *types.Profile, name string) error {
	if err := types.ValidateProfile(Profile); err != nil {
		return err
	}
	var (
		ProjectPath = ""
	)
	if Profile.EnableModules == "on" {
		ProjectPath = path.Join(Profile.GoPath, name)
	} else {
		ProjectPath = path.Join(Profile.GoPath, "src", Profile.VCS, Profile.UserName, name)
	}
	if _, err := os.Stat(ProjectPath); os.IsExist(err) {
		return errors.New("Project folder already exists")
	}
	if err := os.MkdirAll(ProjectPath, os.ModeDir|0751); err != nil {
		return err
	}
	if _, err := git.PlainInit(ProjectPath, false); err != nil {
		return err
	}
	buff := bytes.NewBuffer(nil)
	buff.WriteString(fmt.Sprintf(DefaultReadMeTemplate, strings.ToTitle(name)))
	return ioutil.WriteFile(path.Join(ProjectPath, "README.md"), buff.Bytes(), 0644)
}
