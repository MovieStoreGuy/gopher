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

// CreateProject creates the project based on the profile with the name as the folder
func CreateProject(Profile *types.Profile, name string) error {
	if err := types.ValidateProfile(Profile); err != nil {
		return err
	}
	ProjectPath, err := getProjectRoot(Profile)
	if err != nil {
		return err
	}
	ProjectPath = path.Join(ProjectPath, name)
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
