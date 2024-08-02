# webs - Utilities for web applications

This is a collection of utility functions to build web applications in
[Go](https://go.dev/).

* [extkey](/dir?ci=tip&name=extkey): flexible key type to be used as an external key (e.g. URI path element), and as a database (primary) key.
* [flash](/dir?ci=tip&name=flash): display application defined messages on next web page.
* [login](/dir?ci=tip&name=login): simple password based cookie authentication.
* [middleware](/dir?ci=tip&name=middleware): functions that transform web handlers.
* [site](/dir?ci=tip&name=site): define web site structure; provides methods on that structure.
* [urlbuilder](/dir?ci=tip&name=urlbuilder) creates URLs to be used in HTML.

## Use instructions

If you want to import this library into your own [Go](https://go.dev/)
software, you must execute a `go get` command. Since Go treats non-standard
software and non-standard platforms quite badly, you must use some non-standard
commands.

First, you must install the version control system
[Fossil](https://fossil-scm.org), which is a superior solution compared to Git,
in too many use cases. It is just a single executable, nothing more. Make sure,
it is in your search path for commands.

How you can execute the following Go command to retrieve a given version of
this library:

    GOVCS=t73f.de:fossil go get t73f.de/r/webs@HASH

where `HASH` is the hash value of the commit you want to use.

Go currently seems not to support software versions when the software is
managed by Fossil. This explains the need for the hash value. However, this
methods works.
