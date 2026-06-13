.PHONY: build test install clean release snapshot cross-compile

BINARY  = wlt
VERSION = $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS = -ldflags "-X github.com/weiliantong/cli/internal/build.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o $(BINARY) .

test:
	go test ./...

install: build
	cp $(BINARY) /usr/local/bin/

clean:
	rm -f $(BINARY) $(BINARY).exe
	rm -rf dist
	go clean

# ── Cross-platform release ──────────────────────────────────────────────────
# 构建产物在 dist/，手动上传到云效 CodeUp Release 页面

# GoReleaser snapshot (no tag required, for testing)
snapshot:
	goreleaser release --snapshot --clean
	@echo ""
	@echo "✅ 构建完成，产物在 dist/"
	@echo "   上传到云效 Release: 打开代码库 → 标签 → 发行版 → 新建发行版 → 上传 dist/ 下的文件"
	@echo "   需上传文件:"
	@ls dist/*.tar.gz dist/*.zip dist/checksums.txt 2>/dev/null || true

# GoReleaser full release (requires git tag, 本地构建不自动发布)
release:
	goreleaser release --clean
	@echo ""
	@echo "✅ 构建完成，产物在 dist/"
	@echo "   上传到云效 Release: 打开代码库 → 标签 → 发行版 → 新建发行版 → 上传 dist/ 下的文件"
	@echo "   需上传文件:"
	@ls dist/*.tar.gz dist/*.zip dist/checksums.txt 2>/dev/null || true

# Quick cross-compile without GoReleaser (bare binaries only)
cross-compile:
	@echo "Cross-compiling wlt $(VERSION)..."
	@mkdir -p dist
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/wlt-windows-amd64.exe .
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o dist/wlt-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o dist/wlt-linux-arm64 .
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o dist/wlt-macos-amd64 .
	GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o dist/wlt-macos-arm64 .
	@echo "Done. Binaries in dist/"
