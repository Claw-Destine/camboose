module claw-destine.com/camboose/webapp

go 1.26.2

replace claw-destine.com/camboose/core => ../core

require (
	claw-destine.com/camboose/core v0.0.0-00010101000000-000000000000
	github.com/caarlos0/env/v11 v11.4.1
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/GeertJohan/go.rice v1.0.3 // indirect
	github.com/a-h/parse v0.0.0-20250122154542-74294addb73e // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cli/browser v1.3.0 // indirect
	github.com/cortesi/devd v0.0.0-20260326224400-d0a208f6d905 // indirect
	github.com/cortesi/modd v0.8.1 // indirect
	github.com/cortesi/moddwatch v0.1.0 // indirect
	github.com/cortesi/termlog v0.0.0-20250523085554-f86697764bb0 // indirect
	github.com/daaku/go.zipexe v1.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/goji/httpauth v0.0.0-20160601135302-2da839ab0f4d // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.9.2 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/juju/ratelimit v1.0.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/natefinch/atomic v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rjeczalik/notify v0.9.3 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/toqueteos/webbrowser v1.2.0 // indirect
	golang.org/x/mod v0.35.0 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/term v0.42.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	golang.org/x/tools v0.44.0 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.6.0 // indirect
	gorm.io/gorm v1.31.1 // indirect
	mvdan.cc/sh/v3 v3.6.0 // indirect
)

tool (
	github.com/cortesi/devd/cmd/devd
	github.com/cortesi/modd/cmd/modd
)
