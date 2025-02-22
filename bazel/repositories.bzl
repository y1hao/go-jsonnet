load(
    "@bazel_tools//tools/build_defs/repo:http.bzl",
    "http_archive",
)

# NB: update_cpp_jsonnet.sh looks for these.
CPP_JSONNET_SHA256 = "f104659d3feb42c871d40c5142577dff2cd3b2eda33ac9534f5f12b14643748d"
CPP_JSONNET_GITHASH = "2c3b51491a67ab7b24aecfc23fa3f73f68135129"
CPP_JSONNET_RELEASE_VERSION = "v0.21.0-rc1"

CPP_JSONNET_STRIP_PREFIX = (
    "jsonnet-" + (
        CPP_JSONNET_RELEASE_VERSION if CPP_JSONNET_RELEASE_VERSION else CPP_JSONNET_GITHASH
    )
)
CPP_JSONNET_URL = (
    "https://github.com/google/jsonnet/releases/download/%s/jsonnet-%s.tar.gz" % (
        CPP_JSONNET_RELEASE_VERSION,
        CPP_JSONNET_RELEASE_VERSION,
    ) if CPP_JSONNET_RELEASE_VERSION else "https://github.com/google/jsonnet/archive/%s.tar.gz" % CPP_JSONNET_GITHASH
)

def jsonnet_go_repositories():
    http_archive(
        name = "io_bazel_rules_go",
        sha256 = "33acc4ae0f70502db4b893c9fc1dd7a9bf998c23e7ff2c4517741d4049a976f8",
        urls = [
            "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.48.0/rules_go-v0.48.0.zip",
            "https://github.com/bazelbuild/rules_go/releases/download/v0.48.0/rules_go-v0.48.0.zip",
        ],
    )

    http_archive(
        name = "bazel_gazelle",
        sha256 = "d76bf7a60fd8b050444090dfa2837a4eaf9829e1165618ee35dceca5cbdf58d5",
        urls = [
            "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.37.0/bazel-gazelle-v0.37.0.tar.gz",
            "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.37.0/bazel-gazelle-v0.37.0.tar.gz",
        ],
    )
    http_archive(
        name = "cpp_jsonnet",
        sha256 = CPP_JSONNET_SHA256,
        strip_prefix = CPP_JSONNET_STRIP_PREFIX,
        urls = [CPP_JSONNET_URL],
    )
