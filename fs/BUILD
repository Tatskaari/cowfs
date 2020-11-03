go_library(
    name = "fs",
    srcs = glob(include = ["*.go"], exclude = ["*_test.go"]),
    deps = ["//third_party/go:fuse"],
    visibility = ["PUBLIC"],
)

go_test(
    name = "fs_test",
    srcs = glob(include = ["*_test.go"]),
    deps = [
        ":fs",
        "//third_party/go:testify",
    ],
    data = ["test_data"],
)