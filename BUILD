package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_push")
load("@io_bazel_rules_docker//go:image.bzl", "go_image", GO_DEFAULT_BASE = "DEFAULT_BASE")

proto_library(
    name = "markdown_proto",
    srcs = ["markdown.proto"],
)

go_proto_library(
    name = "markdown_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/jschaf/bazel-go-gce",
    proto = ":markdown_proto",
)

go_binary(
    name = "server",
    srcs = ["main.go"],
    deps = [
        ":markdown_go_proto",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//reflection:go_default_library",
        "@org_golang_x_net//context:go_default_library",
    ],
)

_CONTAINER_PORT = "8282"

# Use a container image so the container is addressable via port-mapping
# on MacOS.
#
# This isn't needed for Linux because --network=host transparently forwards
# all ports between the host and the container.
#
# MacOS uses xhyve to run containers and that apparently makes it impossible
# to publish ports automatically. See below for details:
# https://forums.docker.com/t/should-docker-run-net-host-work/14215/29
container_image(
    name = "go_base_image",
    base = GO_DEFAULT_BASE,
    # MacOS doesn't support --network=host, so work around it by
    # manually publishing the ports:
    # https://github.com/bazelbuild/rules_docker/issues/768
    #
    # --interactive: keeps STDIN open.
    # --rm: removes the container after it stops.
    docker_run_flags = "--interactive --rm --publish=%s:%s" % (_CONTAINER_PORT, _CONTAINER_PORT),
    # Can't specify ports in go_image, so do it in the base image.
    # https://github.com/bazelbuild/rules_docker/issues/309
    ports = [_CONTAINER_PORT],
)

# Must be built with the flag because Docker only runs linux images:
#
# --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64
#
# Once "Flagless multiplatform builds" is supported, we should specify
# Linux AMD64 explicitly in this rule. See the Bazel configuration roadmap:
# https://bazel.build/roadmaps/configuration.html
# https://github.com/bazelbuild/rules_go/issues/2026
go_image(
    name = "server_image",
    base = ":go_base_image",
    binary = ":server",
)

container_push(
    name = "server_push",
    format = "Docker",
    image = ":server_image",
    registry = "gcr.io",
    repository = "experiments-166106/bazel-go-gce",
    tag = "latest",
)
