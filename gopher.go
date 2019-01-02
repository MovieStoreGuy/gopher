package main

import (
	"fmt"
	"os"

	"github.com/MovieStoreGuy/gopher/manager"
	"github.com/MovieStoreGuy/gopher/types"
	"github.com/fatih/color"
	"github.com/voxelbrain/goptions"
)

var (
	options = struct {
		UserProfile string        `goptions:"-p, --profile, description='Define the profile to use with gopher'"`
		Verbose     bool          `goptions:"-v, --verbose, description='Enable verbose logging'"`
		Help        goptions.Help `goptions:"-h, --help, description='Show the global help'"`

		goptions.Verbs
		Create struct {
			Path string `goptions:"-p, --path, description='Define an alternate path to create the project (Requires Go v1.11+)'"`

			Help goptions.Help `goptions:"-h, --help, description='creates a project based on the supplied profile with the name added after flags'"`
			goptions.Remainder
		} `goptions:"create"`
		Project struct {
			Help goptions.Help `goptions:"-h, --help, description='To use projects ensure you supply a nested verb of [show, path, ensure ]'"`
			goptions.Remainder
		} `goptions:"project"`
		Profile struct {
			Name        string `goptions:"-n, --name, description='The name of the profile to store'"`
			VCS         string `goptions:"-v, --vcs, description='The name of the VCS to use'"`
			Username    string `goptions:"-u, --username, description='The username to develop as'"`
			MakeDefault bool   `goptions:"-m, --make-default, description='Make the defined profile the default used'"`

			Help goptions.Help `goptions:"-h, --help, description='To use profile ensure you supply a nested verb of [create, show, set]'"`
			goptions.Remainder
		} `goptions:"profile"`
	}{
		UserProfile: "current",
	}
)

func main() {
	goptions.ParseAndFail(&options)

	currentProfile, err := manager.LoadProfile(options.UserProfile)
	if options.Verbose {
		color.Green("Loaded profile is: %+v", currentProfile)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load profile %s, ensure that the profile exists\n", options.UserProfile)
		fmt.Fprintln(os.Stderr, "Received error: ", err)
		os.Exit(1)
	}
	switch options.Verbs {
	case "create":
		if len(options.Create.Remainder) < 1 {
			color.Yellow("Create required a project name to be supplied")
			os.Exit(1)
		}
		if options.Verbose {
			color.Green("Attempting to create project %s", options.Create.Remainder[0])
		}
		if options.Create.Path != "" {
			currentProfile.GoPath = options.Create.Path
		}
		if manager.CreateProject(currentProfile, options.Create.Remainder[0]); err != nil {
			color.Red("There was an issue trying to create the project")
			color.Red("Error received: %v", err)
			os.Exit(1)
		}
		color.Green("Successfully created project %s", options.Create.Remainder[0])
	case "profile":
		if len(options.Profile.Remainder) < 1 {
			goptions.PrintHelp()
			color.Yellow("Missing subcommand")
			os.Exit(1)
		}
		switch options.Profile.Remainder[0] {
		case "create":
			NewProfile := &types.Profile{
				GoPath:   os.Getenv("GOPATH"),
				Name:     options.Profile.Name,
				VCS:      options.Profile.VCS,
				UserName: options.Profile.Username,
			}
			if options.Verbose {
				color.Green("Creating new profile %s\n", options.Profile.Name)
			}
			if err := manager.StoreProfile(options.Profile.Name, NewProfile); err != nil {
				color.Red("Unable to store profile %s, see error message\n", NewProfile.Name)
				color.Red("Received error: ", err)
				os.Exit(1)
			}
			if options.Profile.MakeDefault {
				if err := manager.SetDefaultProfile(options.Profile.Name); err != nil {
					color.Red("Unable to make default profile %s, see error message\n", NewProfile.Name)
					color.Red("Received error: %v", err)
					os.Exit(1)
				}
			}
			color.Green("Finished creating new profile %s", options.Profile.Name)
		case "set":
			if options.Verbose {
				color.Green("Updating current profile to be %s", options.Profile.Name)
			}
			if err := manager.SetDefaultProfile(options.Profile.Name); err != nil {
				color.Red("Unable to make default profile %s, see error message\n", options.Profile.Name)
				color.Red("Received error: %v", err)
				os.Exit(1)
			}
			color.Green("Successfully updated profile to be %s", options.Profile.Name)
		case "show":
			files, err := manager.GetStoredProfiles()
			if err != nil {
				color.Red("Unable to get stored profiles due to: %v", err)
				os.Exit(1)
			}
			color.Yellow("The stored profiles are:")
			for _, f := range files {
				color.HiWhite("Name:     %s", f.Name)
				color.HiWhite("Username: %s", f.UserName)
				color.HiWhite("VCS:      %s", f.VCS)
				color.HiWhite("---------------")
			}
		default:
			color.Yellow("Unknown option %s", options.Profile.Remainder[0])
			color.Yellow("Please see help for more information")
			return
		}
	case "project":
		if len(options.Project.Remainder) < 1 {
			color.Yellow("Missing required subverbs")
			os.Exit(1)
		}
		switch options.Project.Remainder[0] {
		case "show":
			projects, err := manager.GetProjects(currentProfile)
			if err != nil {
				color.Red("Unable to show projects due to: %v", err)
				os.Exit(1)
			}
			color.Yellow("The current projects stored locally are:")
			for _, p := range projects {
				color.HiWhite("> %s", p)
			}
		case "path":
			if len(options.Project.Remainder) < 2 {
				color.Yellow("path requires a project name")
				os.Exit(1)
			}
			pathName, err := manager.GetProjectPath(currentProfile, options.Project.Remainder[1])
			if err != nil {
				color.Red("Unable to show project due to: %v", err)
				os.Exit(1)
			}
			color.HiWhite("%s", pathName)
		case "ensure":
			var projectName string
			if len(options.Project.Remainder) > 1 {
				projectName = options.Project.Remainder[1]
			}
			if options.Verbose {
				color.Green("Starting ensure process")
			}
			if err := manager.EnsureProjectVendor(currentProfile, projectName); err != nil {
				color.Red("Unable to vendor project")
				color.Red("Received error: %v", err)
				os.Exit(1)
			}
			color.Green("Finished running ensure")
		default:
			color.Yellow("Unknown option %s", options.Project.Remainder[0])
			color.Yellow("Please see help for more information")
			return
		}
	default:
		goptions.PrintHelp()
		color.Yellow("No arguments passed")
	}
}
