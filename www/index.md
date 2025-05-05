# webs - Utilities for web applications

This is a collection of utility functions to build web applications in
[Go](https://go.dev/).

* [aasvg](/dir?ci=tip&name=aasvg): simple ASCII art to SVG translator
* [feed](/dir?ci=tip&name=feed): assist web feed generation.
* [flash](/dir?ci=tip&name=flash): display application defined messages on next web page.
* [ip](/dir?ci=tip&name=ip): working with (remote) addresses in a web context.
* [login](/dir?ci=tip&name=login): simple password based cookie authentication.
* [middleware](/dir?ci=tip&name=middleware): functions that transform web handlers.
* [site](/dir?ci=tip&name=site): define web site structure; provides methods on that structure.
* [urlbuilder](/dir?ci=tip&name=urlbuilder) creates URLs to be used in HTML.
* [xml](/dir?ci=tip&name=feed): assist XML data generation.

## Usage instructions

To import this library into your own [Go](https://go.dev/) software, you need
to run the `go get` command. Since Go does not handle non-standard software and
platforms well, some additional steps are required.

First, install the version control system [Fossil](https://fossil-scm.org),
which is a superior alternative to Git in many use cases. Fossil is just a
single executable, nothing more. Make sure it is included in your system's
command search path.

Then, run the following Go command to retrieve a specific version of
this library:

    GOVCS=t73f.de:fossil go get t73f.de/r/webs@HASH

Here, `HASH` represents the commit hash of the version you want to use.

Go currently does not seem to support software versioning for projects managed
by Fossil. This is why the hash value is required. However, this method works
reliably.
