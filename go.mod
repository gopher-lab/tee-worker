module github.com/masa-finance/tee-worker

go 1.24.0

toolchain go1.24.6

require (
	github.com/edgelesssys/ego v1.8.0
	github.com/google/uuid v1.6.0
	github.com/imperatrona/twitter-scraper v0.0.18
	github.com/joho/godotenv v1.5.1
	github.com/labstack/echo-contrib v0.17.4
	github.com/labstack/echo/v4 v4.13.4
	// FIXME: replace when released
	github.com/masa-finance/tee-types v1.1.18-0.20251009043043-6d5812815e46
	github.com/onsi/ginkgo/v2 v2.26.0
	github.com/onsi/gomega v1.38.2
	github.com/sirupsen/logrus v1.9.3
)

replace github.com/imperatrona/twitter-scraper => github.com/masa-finance/twitter-scraper v1.0.2

require (
	github.com/AlexEidt/Vidio v1.5.1 // indirect
	github.com/Masterminds/semver/v3 v3.4.0 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	go.uber.org/automaxprocs v1.6.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/mod v0.28.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
)

require (
	github.com/go-jose/go-jose/v4 v4.1.3 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20251007162407-5df77e3f7d1d // indirect
	github.com/labstack/gommon v0.4.2
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/exp v0.0.0-20251002181428-27f1f14c8bb9
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	golang.org/x/time v0.14.0 // indirect
	golang.org/x/tools v0.37.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
