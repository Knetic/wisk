#wisk

[![Build Status](https://travis-ci.org/Knetic/wisk.svg?branch=master)](https://travis-ci.org/Knetic/wisk)

`wisk`generates new projects based on a project skeleton. This makes it easy to quickly build new projects without manually creating a lot of boilerplate files and folders.

###Why is this useful?

You (or your company) probably makes projects with very common dependencies, patterns, names, metadata, and boilerplate code for every project. Updating this boilerplate content, or even making a new project, can be daunting, and probably involves a good amount of copy/paste, find/replace, and crossed fingers. With `wisk`, you can make a single skeleton by writing all of the dependencies, structure and conventions once - then re-using that to generate multiple projects.

This tool lets you parameterize file contents, _file names_ and directory names. It seeks to completely eliminate copy/pasting and find/replacing when creating new projects. You should be able to run one command and have a fully-functional (and correctly-named) project ready to go immediately.

###Will this work for my favorite language?

Probably! `wisk` can generate projects for anything that uses text files as its primary mode of representing project structure and data. I personally use it for Java (both standard and maven projects), Go, ruby, python, and Chef.

The only case in which `wisk` will not work is for projects with specific binary-type files involved. If you want `wisk` to work with something like that, you're better off upgrading to a modern technology, or asking your vendor to support modern development workflows.

###How do I use it?

`wisk` takes the path of a skeleton project and the desired output path, then copies the files from the skeleton to the output, like so.

    wisk ./skeleton ./cool_project

If no parameters are given or used, this is equivalent to a `cp` operation. However, `wisk`'s strength is that it can substitute placeholders with parameters. For instance;

    wisk -p "project.name=fooject" ./skeleton ./cool_project

Any placeholders named "project.name" will be replaced by "fooject" in the contents of any file, any file name, or any folder name. So, for instance, in a skeleton file whose contents look like:

    Welcome to ${{=project.name=}}! This project is pretty cool.

After running wisk, the generated file will look like;

    Welcome to fooject! This project is pretty cool.

You can give `wisk` multiple parameters to replace by semicolon-separating them;

    wisk -p "project.name=fooject;project.executable=foo" ./skeleton ./cool_project

Placeholders can be literally anywhere in a plaintext file. `wisk` doesn't parse the file, it just looks for the placeholder tags. You can use them for classnames, module paths, import statements, variable names, README contents, or anywhere else.

If a skeleton contains a parameter that is not specified, the generated project will have all instances of that placeholder replaced with a blank string. This may cause syntax errors, so it's best to always specify every parameter that you need. You can inspect the parameters a skeleton supports by using the `-i` flag, like so;

    wisk -i ./skeleton

This flag works with any valid skeleton, including directories, archives, and registered skeletons.

There are more features to `wisk`, if you're still curious about them you ought to read "ADVANCED.md" in this repository, or check out the "samples" directory - it has examples in many languages, some of which use more advanced features than others. In particular, the "helloworld_ruby" sample shows off every feature in `wisk`.

###How do I make skeletons?

`wisk` accepts any directory as a possible skeleton, just make a folder anywhere and `wisk` will be able to use it. However, this makes it a little hard to share skeletons. So `wisk` also accepts a \*.zip archive. Like so;

    wisk ./skeleton.zip ./cool_project

`wisk` will unzip the file to a temporary directory, then generate a new project based on the contents of the archive.

You can also use "registered templates", which are archives stored per-user. These can be used purely by name, eliminating the need to know exactly where the template is stored. To register a template, create a zip archive of the templat, then use the `-a` flag, like so;

    wisk -a ./skeleton.zip

Afterwards, you can use that same template by name. For example;

    wisk skeleton ./cool_project

Which uses the "skeleton" template (as it was defined when you used the `-a` flag) as the template for "cool_project".

After registration, you can list all registered project skeletons by using the `-l` flag, like so;

    wisk -l

###Can I run a script after wisking a new project?

Absolutely! If the skeleton project contains a file named `_postGenerate.sh` at the top level, then `wisk` will execute that script after generating a project.

The script's working directory will be set to the generated project's directory.

Note that this may have security implications. Inspect post-generation scripts created by others before using project skeletons.

###Aren't there other projects that do this?

Sort of. Maven and Gradle both had the notion of "archetypes", which is conceptually similar. There are projects like Mr. Bones for ruby, and even IDE-specific solutions such as Eclipse or Visual Studio's project templates. That's to say nothing of all the framework-specific project generators. But I've grown tired of learning a new template scheme and cli interface for every language or framework, I just wanted to have one set of templates that were all user-defined, easy to share, and all using a unified parameter interface.

[Cookiecutter](https://github.com/audreyr/cookiecutter) is a cool project which meets similar goals, but has many drawbacks. It requires a python installation, needs either STDIN input or an *rc file, has logic-based templating, can't create templates from arbitrary directories or archives, seems to aim for monolithic templates, cannot do parameter joining or recursion, and doesn't seem to support post-generation scripts; making it difficult (if not impossible) to use for languages like Ruby, C-sharp, or Java (which I, at least, spend a lot of time in). `wisk` uses small, direct, purposeful skeletons to create many similar projects, it's not meant to make a monolithic language/framework template which fills every conceivable need that one may have of a project in that language/framework.

###Branching

I use green masters, and heavily develop with private feature branches. Full releases are pinned and unchangeable, representing the best available version with the best documentation and test coverage. Master branch, however, should always have all tests pass and implementations considered "working", even if it's just a first pass. Master should never panic.

###License

This project is licensed under the MIT general use license. You're free to integrate, fork, and play with this code as you feel fit without consulting the author, as long as you provide proper credit to the author in your works.

###Activity

If this repository hasn't been updated in a while, it's probably because I don't have any outstanding issues to work on - it's not because I've abandoned the project. If you have questions, issues, or patches; I'm completely open to pull requests, issues opened on github, or emails from out of the blue.
