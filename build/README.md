# `/build`

Building, Packaging and Continuous Integration info.

## Building

When some libraries without stable versions require particular Go version which is not yet provided by the Red Hat go-toolset, we need to switch to building with upstream version of Go until the particular version is available. To do this, switch the `Dockerfile` symlink to one of the following build configurations:

* `RedHat.dockerfile` - builds the project via Go compiler from the Red Hat go-toolset image.
* `Upstream.dockerfile` - builds the project via Go compiler from the Go upstream project.

## Packaging

Put your cloud (AMI), container (Docker), OS (deb, rpm, pkg) package configurations and scripts in the `/build/package` directory.

## CI

Put your CI (travis, circle, drone) configurations and scripts in the `/build/ci` directory. Note that some of the CI tools (e.g., Travis CI) are very picky about the location of their config files. Try putting the config files in the `/build/ci` directory linking them to the location where the CI tools expect them when possible (don't worry if it's not and if keeping those files in the root directory makes your life easier :-)).

Examples:

* https://github.com/cockroachdb/cockroach/tree/master/build

## Updating dependencies

To update all our project dependencies to the latest version:

    go get -u ./... && go mod tidy
