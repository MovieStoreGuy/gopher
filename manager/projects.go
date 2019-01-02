package manager

import (
	"github.com/MovieStoreGuy/gopher/types"
	"github.com/fatih/color"
	"github.com/golang/dep"
	"gopkg.in/src-d/go-git.v4"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
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

func EnsureProjectVendor(Profile *types.Profile, name string) error {
	if err := types.ValidateProfile(Profile); err != nil {
		return err
	}
	var projectDirectory string
	if name != "" {
		root, err := getProjectRoot(Profile)
		if err != nil {
			return err
		}
		projectDirectory = path.Join(root, name)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		projectDirectory = wd
	}

	ctx := &dep.Ctx{}
	GOPATHS := filepath.SplitList(Profile.GoPath)
	ctx.SetPaths(projectDirectory, GOPATHS...)

	project, err := ctx.LoadProject()
	if err != nil {
		return err
	}

	sourcemanager, err := ctx.SourceManager()
	if err != nil {
		return err
	}

	sourcemanager.UseDefaultSignalHandling()
	defer sourcemanager.Release()

	if err := dep.ValidateProjectRoots(ctx, project.Manifest, sourcemanager); err != nil {
		return err
	}

	packages, err := project.GetDirectDependencyNames(sourcemanager)
	if err != nil {
		return err
	}

	checkpath, err := getProjectRoot(Profile)
	if err != nil {
		return err
	}
	checkpath = strings.Split(checkpath, "src/")[1]
	sshauth, err := gitSShAuthConfig(Profile)
	if err != nil {
		return err
	}
	for importpackage := range packages {
		p := string(importpackage)
		color.HiWhite("Project: %s, checkout Path: %s",p, checkpath)
		if strings.HasPrefix(p, checkpath) {
			os.MkdirAll(path.Dir(path.Join("vendor", p)), os.ModeDir|0751)
			color.HiWhite("Clone path %s", "git@" + strings.Replace(p, "/", ":", 1) + ".git")
			_, err := git.PlainClone(path.Join("vendor", p), true, &git.CloneOptions{
				URL:  "git@" + strings.Replace(p, "/", ":", 1) + ".git",
				Auth: &sshauth,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
