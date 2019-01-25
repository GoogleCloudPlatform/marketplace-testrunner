# Testrunner

Testrunner is a binary and docker container for running common tests
without requiring to write code. It's a simple and fast way to add
functional testing to an application, as its test spec is defined in
a `yaml` file.

Testrunner also has useful integrations with GCP Deployment Manager and
Kubernetes applications.

# Development

## Preprequisites

This repository uses [bazel](https://bazel.build) to build the binary and
container. It also supports [cloudbuild](https://cloud.google.com/container-builder/docs/)
to build and publish your container on GCP from source.

The repository is also compatible with `go` tool. You'll need to
install dependencies separately as they are not vendored.

## Generate BUILD files

```shell
bazel run //:gazelle
```

## Build locally


### Binary

Build and run the binary:

```shell
bazel run //runner:main -- -logtostderr --test_spec=$PWD/examples/testspecs/http.yaml
```

### Container

To build and run the docker container:

```shell
# Run tests
bazel test //...

# Build binary
bazel build //runner:main

# Make temporary directory
mkdir -p tmp

# Copy the file and rename it
cp bazel-bin/runner/testrunner tmp/testrunner

# Copy all Docker specific files
cp docker/* tmp

# Build container
docker build --tag=testrunner tmp

# Run the installed container, mounting the test definition
# files as a volume.
docker run --rm \
  -v=$PWD/examples:/examples \
  testrunner -logtostderr --test_spec=/examples/testspecs/http.yaml
```

## Build GCP container

Two workarounds before running the command below:

Execute the following.

```shell
rm -r bazel-*
```

Then execute the following.

```shell
gcloud builds submit --config cloudbuild.yaml .
```

This publishes a `testrunner` container in your project (i.e.
whatever the default project for your `gcloud` is).

You can test by pulling the published image and run it.

```shell
export PROJECT=$(gcloud config get-value project)

docker run --rm \
  -v=$PWD/examples:/examples \
  gcr.io/$PROJECT/testrunner -logtostderr --test_spec=/examples/testspecs/http.yaml
```
