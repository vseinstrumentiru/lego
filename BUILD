github_repo(
    name = "pleasings2",
    repo = "sagikazarmark/mypleasings",
    revision = "master",
)

subinclude("///pleasings2//go")

go_build(
    name = "lego",
    labels = ["binary"],
    stamp = True,
    trimpath = True,
)

tarball(
    name = "artifact",
    srcs = ["README.md", ":lego"],
    out = f"lego_{CONFIG.OS}_{CONFIG.ARCH}.tar.gz",
    gzip = True,
    labels = ["manual"],
)

subinclude("///pleasings2//misc")

build_artifacts(
    name = "artifacts",
    artifacts = [
        "@linux_amd64//:artifact",
        "@darwin_amd64//:artifact",
    ],
    labels = ["manual"],
)

subinclude("///pleasings2//github")

github_release(
    name = "publish",
    assets = [":artifacts"],
    labels = ["manual"],
)
