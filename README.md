# A.P.O: - Authorized personnel only

APO is a simple OAuth IAM server.
It's backed is simple htpasswd file (as used e.g. in Apache)
It issues tokens under `/token`.
Tokens are signed using ECDSA, using a key generated at startup and exposed as jwk under `/key`

**Please Note:** this app server currently only provides an http endpoint (not https),
so you should expose it behind a reverse proxy (such as nginx).
This allows you also to use a certificate managed by letsencrypt.

## Building, staring

`go run main.go`

## Configuration

config is based on viper. The server looks for a file called `iam-conf.json` and evaluates command line parameters and environment variables 

### passfile

set via command line using `-f` or `--passfile` resp `"passfile"`in the config file.

This tells the server which htpasswd file to use as user base. defaults to ``"./passes"`` 
Not re-evaluated at runtime, restart IAM on changes.

### endpoint

set via command line using `-e` or `--endpoint` resp `"endpoint"`in the config file.

where this server listens. defaults to ``":3000"``

### audience

set via command line using `-a` or `--audience` resp `"audience"`in the config file.

the audience/realm that gets protected. Appears as such in the basic auth challenge of IAM and in the issued token.
defaults to `"myhome"`

## dependencies

### direct
```
	"github.com/Sirupsen/logrus"
	"github.com/abbot/go-http-auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/foomo/htpasswd"
	"github.com/gorilla/mux"
	"github.com/mendsley/gojwk"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
```

### transitive
```
github.com/Sirupsen/logrus          MIT License
github.com/abbot/go-http-auth       Apache License 2.0 (95%)
github.com/dgrijalva/jwt-go         MIT License (98%)
github.com/foomo/htpasswd           MIT License
github.com/fsnotify/fsnotify        BSD 3-clause "New" or "Revised" License (96%)
github.com/gorilla/mux              BSD 3-clause "New" or "Revised" License (96%)
github.com/h0ru5/hmauth/IAM         MIT License (98%)
github.com/hashicorp/hcl            Mozilla Public License 2.0
github.com/kr/fs                    BSD 3-clause "New" or "Revised" License (96%)
github.com/magiconair/properties    BSD 2-clause "Simplified" License (95%)
github.com/mendsley/gojwk           BSD 2-clause "Simplified" License (96%)
github.com/mitchellh/mapstructure   MIT License
github.com/pelletier/go-buffruneio  ?
github.com/pelletier/go-toml        MIT License
github.com/pkg/errors               BSD 2-clause "Simplified" License
github.com/pkg/sftp                 BSD 2-clause "Simplified" License
github.com/spf13/afero              Apache License 2.0 (95%)
github.com/spf13/cast               MIT License
github.com/spf13/jwalterweatherman  MIT License
github.com/spf13/pflag              BSD 3-clause "New" or "Revised" License (96%)
github.com/spf13/viper              MIT License
golang.org/x/crypto                 BSD 3-clause "New" or "Revised" License (96%)
golang.org/x/net/context            BSD 3-clause "New" or "Revised" License (96%)
golang.org/x/sys/unix               BSD 3-clause "New" or "Revised" License (96%)
golang.org/x/text                   BSD 3-clause "New" or "Revised" License (96%)
gopkg.in/yaml.v2                    ? (The Unlicense, 35%)
```

## Licence / Copyright

Copyright (c) Johannes Hund / @h0ru5 2017

Published under the [MIT licence](https://opensource.org/licenses/MIT)
