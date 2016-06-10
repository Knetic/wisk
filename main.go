package main

import (
	"fmt"
	"os"
)

func main() {

	var registry *TemplateRegistry
	var settings RunSettings
	var err error

	registry = NewTemplateRegistry()

	settings, err = FindRunSettings()
	if err != nil {
		exitWith("Unable to parse run arguments: %s\n", err, 1)
		return
	}

	if(settings.flagList) {

		for _, flag := range FLAGS {
			fmt.Println(flag)
		}
		return
	}

	registry = NewTemplateRegistry()

	// is the user showing the registry?
	if settings.showRegistry {
		showRegistry(registry)
		return
	}

	// is the user trying to add to the current registry?
	if settings.addRegistry {
		addRegistry(settings, registry)
		return
	}

	createProject(settings, registry)
}

func showRegistry(registry *TemplateRegistry) {

	for _, template := range registry.GetTemplateList() {
		fmt.Println(template)
	}
}

func addRegistry(settings RunSettings, registry *TemplateRegistry) {

	var name string
	var err error

	name, err = registry.RegisterTemplate(settings.skeletonPath, settings.basicAuthUser, settings.basicAuthPass)
	if err != nil {

		// TODO: I'm deeply uncomfortable with using "exitWith" outside of the actual
		// main method. This is too easy to let "exiting" become a separate code path.
		exitWith("Unable to register template: %s\n", err, 1)
	}

	fmt.Printf("Registered template '%s'\n", name)
	return
}

func createProject(settings RunSettings, registry *TemplateRegistry) {

	var project *TemplatedProject
	var parameters []string
	var err error

	// is this a registry skeleton?
	if registry.IsPathRegistry(settings.skeletonPath) && registry.Contains(settings.skeletonPath) {

		settings.skeletonPath, err = registry.GetTemplatePath(settings.skeletonPath)
		if err != nil {
			exitWith("Unable to expand registered template: %s\n", err, 1)
			return
		}
	}

	// Create templated project, in preparation for later use
	project, err = NewTemplatedProject(settings.skeletonPath)
	if err != nil {
		exitWith("Unable to read templated project: %s\n", err, 1)
		return
	}

	// inspect only?
	if settings.inspectionRun {

		parameters, err = project.FindParameters()

		if err != nil {
			exitWith("Unable to inspect skeleton: %s\n", err, 1)
			return
		}

		for _, parameter := range parameters {
			fmt.Println(parameter)
		}
		return
	}

	// check to see if the destination exists
	_, err = os.Stat(settings.targetPath)
	if(err == nil) {

		// if we force-generate, see which type of force generation to do (delete or overwrite)
		switch settings.forceGenerate {
		case GENERATE_NONE:
			fmt.Println("Destination path already exists, and no '-f' option was specified. Use '-f' to overwrite existing files.")
			return
		case GENERATE_DELETE:
			err = os.RemoveAll(settings.targetPath)
			if(err != nil) {
				exitWith("Unable to delete existing files in the given output directory: %s\n", err, 1)
				return
			}
		default:
			break
		}
	}

	// generate a project
	err = project.GenerateAt(settings.targetPath, settings.parameters)
	if err != nil {
		exitWith("Unable to generate project: %s\n", err, 1)
		return
	}

	// if everything succeeded, but we had missing parameters, make a note of it to the user.
	if project.missingParameters.Length() > 0 {

		if(settings.blankParameters) {

			// if blank parameters are allowed, print a message about how blank parameters were used.
			fmt.Printf("Project generated, but some parameters were not specified, and have been left blank:\n")
		} else {

			// if blank parameters are not allowed, delete the generated project.
			fmt.Printf("Project cannot be generated, some parameters were not specified:\n")
			os.RemoveAll(settings.targetPath)
		}

		for _, value := range project.missingParameters.GetSlice() {
			fmt.Println(value)
		}
		return
	}

	if project.incorrectParameters.Length() > 0 {

		fmt.Println("Project generated, but some placeholders in the template were specified with a join character ('[]'), but the given parameter was not a list")
		fmt.Println("This may mean you need to specify the following parameters as a comma-separated list")

		for _, value := range project.incorrectParameters.GetSlice() {
			fmt.Println(value)
		}
	}

	// if there's a post-generate script (and it's executable), call it.
	err = executePostGenerate(project.rootDirectory, settings.targetPath)
	if err != nil {
		exitWith("Unable to complete post-generation script: %s\n", err, 1)
		return
	}
}

func exitWith(message string, err error, code int) {

	errorMsg := fmt.Sprintf(message, err)
	fmt.Fprintf(os.Stderr, errorMsg)
	os.Exit(code)
}
