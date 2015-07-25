#wisk


`wisk`generates new projects based on a project skeleton. This makes it easy to quickly build new projects without manually creating a lot of boilerplate files and folders.

###Why is this useful?


I tend to make _lots_ of small projects. Utilities, gems, modules, websites, services - everything. And every language and framework requires a different project file structure, meaning that in order to even build "hello world" you need to bear a lot of project-specific file structures and contents in mind. `wisk` allows you to build a skeleton _once_, then reuse it to speed up the creation of new projects.

`wisk` makes more of a difference with more complex templates. Your company probably makes projects with very common dependencies, patterns, names, and boilerplate code for every project. Updating this boilerplate content, or even making a new project, can be daunting without using a skeleton project. With `wisk`, you  can make a single skeleton with all of the dependencies, structure, conventions, and written once, and used to generate many projects.

###Will this work for my favorite language?

Probably! `wisk` can generate projects for anything that uses text files as its primary mode of representing project structure and data. I personally use it for Java (both standard and maven projects), Go, ruby, python, and Chef.

The only case in which `wisk` will not work is for project with specific proprietary mechanisms for creating projects; such as anything .NET-related. If you want `wisk` to work with something like that, you're better off upgrading to a modern technology, or asking your vendor to support modern development workflows.

###How do I use it?

`wisk` takes the path of a skeleton project, and the desired output path, and copies the files from the skeleton to the output, like so.

    wisk ./skeleton ./cool_project

If no parameters are given or used, this is equivalent to a `cp` operation. However, `wisk`'s strength is that it can substitute placeholders with parameters. For instance;

    wisk -p "project.name=fooject" ./skeleton ./cool_project

Any placeholders named "project.name" will be replaced by "fooject" in the contents of any file, any file name, or any folder name. So, for instance, in a skeleton file whose contents look like:

    Welcome to ${{=project.name=}}! This project is pretty cool.

After running wisk, the generated file will look like;

    Welcome to fooject! This project is pretty cool.

You can give `wisk` multiple parameters to replace by semicolon-separating them;

    wisk -p "project.name=fooject;project.executable=foo" ./skeleton ./cool_project

If a skeleton contains a parameter that is not specified, a warning is printed informing you of that. In that case, the generated project will have all instances of that placeholder replaced with a blank string. This may cause syntax errors, so it's best to always specify every parameter that you need.

###How do I make skeletons?

`wisk` accepts any path as a possible skeleton, just make a folder anywhere and `wisk` will use it. However, this makes it a little hard to share skeletons. So `wisk` also accepts a \*.zip archive. Like so;

    wisk ./skeleton.zip ./cool_project

`wisk` will unzip the file to a temporary directory, then generate a new project based on the contents of the archive.

You can also use "registered templates", which are archives stored per-user. These can be used purely by name, eliminating the need to know exactly where the template is stored. To register a template, create a zip archive of the templat, then use the `-a` flag, like so;

    wisk -a ./skeleton.zip

Afterwards, you can use that same template by name. For example;

    wisk skeleton ./cool_project

Which uses the "skeleton" template (as it was defined when you used the `-a` flag) as the template for "cool_project".

After registration, you can list all registered project skeletons by using the `-l` flag, like so;

    wisk -l

###How do I know what parameters a skeleton accepts?

Running `wisk` with the "-i" flag will inspect the given skeleton, and print out a list of all parameters used by it. Like so;

    wisk -i ./skeleton

This flag works with any valid skeleton, including directories, archives, and registered skeletons. So all of the following work the same (assuming they all exist);

    wisk -i ./skeleton
    wisk -i ./skeleton.zip
    wisk -i skeleton

###Can I run a script after wisking a new project?

Absolutely! If the skeleton project contains a file named `_postGenerate.sh` at the top level, then `wisk` will execute that script after generating a project.

The script's working directory will be set to the generated project's directory.

Note that this may have security implications. Inspect post-generation scripts created by others before using project skeletons.

###Branching

I use green masters, and heavily develop with private feature branches. Full releases are pinned and unchangeable, representing the best available version with the best documentation and test coverage. Master branch, however, should always have all tests pass and implementations considered "working", even if it's just a first pass. Master should never panic.

###License

This project is licensed under the MIT general use license. You're free to integrate, fork, and play with this code as you feel fit without consulting the author, as long as you provide proper credit to the author in your works.

###Activity

If this repository hasn't been updated in a while, it's probably because I don't have any outstanding issues to work on - it's not because I've abandoned the project. If you have questions, issues, or patches; I'm completely open to pull requests, issues opened on github, or emails from out of the blue.
