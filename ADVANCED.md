# Advanced usage

###How do I make package paths?

Parameters can be specified as a list by seprating values with a comma, like so.

    wisk -p "project.package=com,example,sample"

By default, if a parameter is referenced and the parameter is a list, only the first element in the list
will be used. So using "${{=project.package=}}" with the above parameters will only result in "com".
However, placeholders can specify a separator used to join the list elements together, like so;

    ${{=project.package[.]=}}

In this case, the above will be replaced with "com.example.sample". This is useful for creating nested folder structures, or package names. See the Java examples in the "samples" folder for an implementation of this, using a single "project.package" parameter to create nested folder hierarchies and package declarations.

Note that if no separator is specified, the default OS path separator is used instead.

###Can I inject parameters into a block of boilerplate?

A "content placeholder" can be used to create a sequence of values, each using one value from the list of a single parameter. For instance;

	${{=:project.properties=}}
	${{=_value=}}=TRUE
	${{=;project.properties=}}

Would use each value in "project.properties" to generate a line. Given the input:

	wisk -p "project.properties=foo,bar,baz,quux"

The following would be generated:

	foo=TRUE
	bar=TRUE
	baz=TRUE
	quux=TRUE

However, this "content placeholder" construct can be used recursively with the `_recurse_` reserved placeholder. This is primarily useful for things like Ruby module declarations. Given the below example;

	${{=:project.module=}}
	module ${{=_value=}}
		${{=_recurse=}}
	end
	${{=;project.module=}}

With the parameter being;

	wisk -p "project.module=My,Sample,Project"

The result is:

	module My
		module Sample
			module Project
			end
		end
	end

Note that as of this time content placeholders cannot be nested. The following is _invalid_:

	${{=:project.module=}}
	some
	content
	${{=_value=}}
		${{=:project.innerModule=}}
		more
		content
		${{=_value=}}
		${{=;project.innerModule=}}
	${{=;project.module=}}

There will be an issue opened to support this, but it won't be on the roadmap unless some use cases make themselves apparent. At present, it doesn't seem too useful for any projects I can think of. But if you want this functionality, +1 the issue (or open a new one).

###Will this overwrite existing files?

Yes. `wisk` will overwrite any file completely - it has no merging strategy, and does not warn you that it is going to do this. This is intentional - you may wish to use `wisk` incrementally, regenerating boilerplate project files multiple times as you update the skeleton.
