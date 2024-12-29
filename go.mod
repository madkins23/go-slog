module github.com/madkins23/go-slog

go 1.22.0

toolchain go1.22.5

require (
	github.com/dmarkham/enumer v1.5.10
	github.com/fatih/camelcase v1.0.0
	github.com/gertd/go-pluralize v0.2.1
	github.com/gin-gonic/gin v1.10.0
	github.com/gomarkdown/markdown v0.0.0-20241205020045-f7e15b2f3e62
	github.com/madkins23/gin-utils v1.4.1
	github.com/madkins23/go-utils v1.44.0
	github.com/phsym/console-slog v0.3.1
	github.com/phsym/zeroslog v0.2.0
	github.com/phuslu/log v1.0.113
	github.com/rs/zerolog v1.33.0
	github.com/samber/slog-logrus/v2 v2.5.1
	github.com/samber/slog-zap/v2 v2.6.1
	github.com/samber/slog-zerolog/v2 v2.7.2
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.10.0
	github.com/veqryn/slog-dedup v0.5.0
	github.com/vicanso/go-charts/v2 v2.6.10
	github.com/wcharczuk/go-chart/v2 v2.1.1
	go.mrchanchal.com/zaphandler v0.0.0-20230611140024-bd4fd80897ad
	go.seankhliao.com/svcrunner/v3 v3.0.0-20231007180458-c5294d90b36c
	go.uber.org/zap v1.27.0
	golang.org/x/text v0.21.0
	snqk.dev/slog/meld v0.0.0-20240701183407-595424398869
)

// This breaks phsym/zeroslog
// TODO: Remove? exclude github.com/rs/zerolog v1.32.0

// Changed
// - c.w.Write([]byte(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="%d" height="%d">\n`, c.width, c.height)))
// + c.w.Write([]byte(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 %d %d">`, c.width, c.height)))
// which broke size in Chrome on Linux.
// Older version still not correct, should be BOTH viewBox and width/height.
exclude github.com/wcharczuk/go-chart/v2 v2.1.2

// TODO: Remove this when phsym/zeroslog merges PR#6
// TODO: Remove? exclude github.com/rs/zerolog v1.33.0

require (
	github.com/blend/go-sdk v1.20240719.1 // indirect
	github.com/bytedance/sonic v1.12.6 // indirect
	github.com/bytedance/sonic/loader v0.2.1 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.7 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.23.0 // indirect
	github.com/goccy/go-json v0.10.4 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.9 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pascaldekloe/name v1.0.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/samber/lo v1.47.0 // indirect
	github.com/samber/slog-common v0.18.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	go.opentelemetry.io/otel v1.33.0 // indirect
	go.opentelemetry.io/otel/trace v1.33.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.12.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/image v0.23.0 // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/tools v0.28.0 // indirect
	google.golang.org/protobuf v1.36.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/b/v2 v2.1.1 // indirect
)
