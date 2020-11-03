package(default_visibility = ["PUBLIC"])

go_get(
    name = "fuse",
    get = "bazil.org/fuse",
    install = [
        ".",
        "fs",
        "fuseutil",
    ],
    deps = [":xsys"],
)

go_get(
    name = "xsys",
    get = "golang.org/x/sys/...",
)

go_get(
    name = "go-cli-init",
    get = "github.com/peterebden/go-cli-init",
    revision = "v1.3.1",
    deps = [
        ":go-flags",
        ":humanize",
        ":logging",
        ":terminal",
    ],
)

go_get(
    name = "go-flags",
    get = "github.com/jessevdk/go-flags",
    revision = "v1.4.0",
)

go_get(
    name = "humanize",
    get = "github.com/dustin/go-humanize",
    revision = "8929fe90cee4b2cb9deb468b51fb34eba64d1bf0",
)

go_get(
    name = "logging",
    get = "gopkg.in/op/go-logging.v1",
    revision = "b2cb9fa56473e98db8caba80237377e83fe44db5",
)

go_get(
    name = "terminal",
    get = "golang.org/x/crypto/ssh/terminal",
    revision = "7b85b097bf7527677d54d3220065e966a0e3b613",
)

go_get(
    name = "testify",
    get = "github.com/stretchr/testify",
    install = [
        "assert",
        "require",
        "suite",
        "vendor/...",
    ],
    licences = ["MIT"],
    revision = "v1.4.0",
    deps = [
        ":difflib",
        ":objx",
        ":spew",
        ":yaml.v2",
    ],
)

go_get(
    name = "difflib",
    get = "github.com/pmezard/go-difflib/...",
    licences = ["BSD 3-Clause"],
    revision = "792786c7400a136282c1664665ae0a8db921c6c2",
)

go_get(
    name = "objx",
    get = "github.com/stretchr/objx",
    licences = ["MIT"],
    revision = "1a9d0bb9f541897e62256577b352fdbc1fb4fd94",
)

go_get(
    name = "yaml.v2",
    get = "gopkg.in/yaml.v2",
    licences = ["Apache 2.0"],
    revision = "v2.2.4",
)

go_get(
    name = "spew",
    get = "github.com/davecgh/go-spew/spew",
    licences = ["ISC"],
)