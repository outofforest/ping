package build

import (
	"context"
	"os/exec"

	"github.com/outofforest/build/v2/pkg/tools"
	"github.com/outofforest/build/v2/pkg/types"
	"github.com/outofforest/libexec"
	"github.com/outofforest/tools/pkg/tools/golang"
)

// sudo setcap cap_net_raw=eip bin/ping-app

func buildApp(ctx context.Context, deps types.DepsFunc) error {
	deps(golang.EnsureGo)

	return golang.Build(ctx, deps, golang.BuildConfig{
		Platform:      tools.PlatformLocal,
		PackagePath:   ".",
		BinOutputPath: "bin/ping-app",
	})
}

func runApp(ctx context.Context, deps types.DepsFunc) error {
	deps(buildApp)

	return libexec.Exec(ctx, exec.Command("bin/ping-app"))
}
