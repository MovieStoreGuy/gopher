package manager

import (
	"github.com/MovieStoreGuy/gopher/types"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// GetProjects gets all projects under the current profile
func GetProjects(Profile *types.Profile) ([]string, error) {
	projectRoot, err := getProjectRoot(Profile)
	if err != nil {
		return nil, err
	}
	files, err := ioutil.ReadDir(projectRoot)
	if err != nil {
		return nil, err
	}
	var projects []string
	for _, f := range files {
		if !strings.HasPrefix(f.Name(), ".") && f.IsDir() {
			if _, err := os.Stat(path.Join(projectRoot, f.Name(), ".git")); !os.IsNotExist(err) {
				projects = append(projects, f.Name())
			} else {
				for _, project := range getNestedProjects(path.Join(projectRoot, f.Name())) {
					projects = append(projects, path.Join(f.Name(), project))
				}
			}
		}
	}
	return projects, nil
}

func getNestedProjects(project string) []string {
	var projects []string
	files, err := ioutil.ReadDir(project)
	if err != nil {
		return projects
	}
	for _, file := range files {
		if file.IsDir() {
			if _, err := os.Stat(path.Join(project, file.Name(), ".git")); !os.IsNotExist(err) {
				projects = append(projects, file.Name())
				continue
			}
			nestedProjects := getNestedProjects(path.Join(project, file.Name()))
			for _, nested := range nestedProjects {
				projects = append(projects, path.Join(file.Name(), nested))
			}
		}
	}
	return projects
}

// GetProjectPath builds the profile path based of the Profile and returns to back to the callee
func GetProjectPath(Profile *types.Profile, name string) (string, error) {
	if err := types.ValidateProfile(Profile); err != nil {
		return "", err
	}
	projectRoot, err := getProjectRoot(Profile)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(path.Join(projectRoot, name)); os.IsNotExist(err) {
		return "", err
	}
	return path.Join(projectRoot, name), nil
}
