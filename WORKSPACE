# Go build setup.
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "67b4d1f517ba73e0a92eb2f57d821f2ddc21f5bc2bd7a231573f11bd8758192e",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.50.0/rules_go-v0.50.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.50.0/rules_go-v0.50.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "b760f7fe75173886007f7c2e616a21241208f3d90e8657dc65d36a771e916b6a",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.39.1/bazel-gazelle-v0.39.1.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.39.1/bazel-gazelle-v0.39.1.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.23.0")

gazelle_dependencies()

load("@bazel_gazelle//:deps.bzl", "go_repository")

go_repository(
    name = "in_gopkg_yaml_v2",
    importpath = "gopkg.in/yaml.v2",
    tag = "v2",
)

go_repository(
    name = "in_gopkg_xmlpath_v2",
    importpath = "gopkg.in/xmlpath.v2",
    tag = "v2",
)

go_repository(
    name = "org_golang_x_net",
    importpath = "golang.org/x/net",
    tag = "release-branch.go1.8",
)

go_repository(
    name = "com_github_golang_glog",
    importpath = "github.com/golang/glog",
    tag = "master",
)

go_repository(
    name = "com_github_ghodss_yaml",
    importpath = "github.com/ghodss/yaml",
    tag = "master",
)

