package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	chserver "github.com/zgsm-ai/cotun/server"
	chshare "github.com/zgsm-ai/cotun/share"
	"github.com/zgsm-ai/cotun/share/ccrypto"
	"github.com/zgsm-ai/cotun/share/cos"
	"github.com/zgsm-ai/cotun/share/settings"
)

var commonHelp = `
    --pid Generate pid file in current working directory

    -v, Enable verbose logging

    --help, This help text

  Signals:
    The cotun process is listening for:
      a SIGUSR2 to print process stats, and
      a SIGHUP to short-circuit the client reconnect timer

  Version:
    ` + chshare.BuildVersion + ` (` + runtime.Version() + `)

  Read more:
    https://github.com/zgsm-ai/cotun

`

func generatePidFile() {
	pid := []byte(strconv.Itoa(os.Getpid()))
	if err := os.WriteFile("cotun.pid", pid, 0644); err != nil {
		log.Fatal(err)
	}
}

var serverHelp = `
  Usage: cotund [options]

  Options:

    --host, Defines the HTTP listening host – the network interface
    (defaults the environment variable HOST and falls back to 0.0.0.0).

    --port, -p, Defines the HTTP listening port (defaults to the environment
    variable PORT and fallsback to port 8080).

    --key, (deprecated use --keygen and --keyfile instead)
    An optional string to seed the generation of a ECDSA public
    and private key pair. All communications will be secured using this
    key pair. Share the subsequent fingerprint with clients to enable detection
    of man-in-the-middle attacks (defaults to the COTUN_KEY environment
    variable, otherwise a new key is generate each run).

    --keygen, A path to write a newly generated PEM-encoded SSH private key file.
    If users depend on your --key fingerprint, you may also include your --key to
    output your existing key. Use - (dash) to output the generated key to stdout.

    --keyfile, An optional path to a PEM-encoded SSH private key. When
    this flag is set, the --key option is ignored, and the provided private key
    is used to secure all communications. (defaults to the COTUN_KEY_FILE
    environment variable). Since ECDSA keys are short, you may also set keyfile
    to an inline base64 private key (e.g. cotun server --keygen - | base64).

    --authfile, An optional path to a users.json file. This file should
    be an object with users defined like:
      {
        "<user:pass>": ["<addr-regex>","<addr-regex>"]
      }
    when <user> connects, their <pass> will be verified and then
    each of the remote addresses will be compared against the list
    of address regular expressions for a match. Addresses will
    always come in the form "<remote-host>:<remote-port>" for normal remotes
    and "R:<local-interface>:<local-port>" for reverse port forwarding
    remotes. This file will be automatically reloaded on change.

    --auth, An optional string representing a single user with full
    access, in the form of <user:pass>. It is equivalent to creating an
    authfile with {"<user:pass>": [""]}. If unset, it will use the
    environment variable AUTH.

    --keepalive, An optional keepalive interval. Since the underlying
    transport is HTTP, in many instances we'll be traversing through
    proxies, often these proxies will close idle connections. You must
    specify a time with a unit, for example '5s' or '2m'. Defaults
    to '25s' (set to 0s to disable).

    --backend, Specifies another HTTP server to proxy requests to when
    cotun receives a normal HTTP request. Useful for hiding cotun in
    plain sight.

    --socks5, Allow clients to access the internal SOCKS5 proxy. See
    cotun client --help for more information.

    --reverse, Allow clients to specify reverse port forwarding remotes
    in addition to normal remotes.

    --tls-key, Enables TLS and provides optional path to a PEM-encoded
    TLS private key. When this flag is set, you must also set --tls-cert,
    and you cannot set --tls-domain.

    --tls-cert, Enables TLS and provides optional path to a PEM-encoded
    TLS certificate. When this flag is set, you must also set --tls-key,
    and you cannot set --tls-domain.

    --tls-domain, Enables TLS and automatically acquires a TLS key and
    certificate using LetsEncrypt. Setting --tls-domain requires port 443.
    You may specify multiple --tls-domain flags to serve multiple domains.
    The resulting files are cached in the "$HOME/.cache/cotun" directory.
    You can modify this path by setting the COTUN_LE_CACHE variable,
    or disable caching by setting this variable to "-". You can optionally
    provide a certificate notification email by setting COTUN_LE_EMAIL.

    --tls-ca, a path to a PEM encoded CA certificate bundle or a directory
    holding multiple PEM encode CA certificate bundle files, which is used to
    validate client connections. The provided CA certificates will be used
    instead of the system roots. This is commonly used to implement mutual-TLS.

    --minport, Defines the minimum port for port allocation (defaults to 30000).

    --maxport, Defines the maximum port for port allocation (defaults to 31000).

    --control-port, Control plane port, used for managing port information. Supports the following API endpoints:
      GET /{moduleName}/api/v1/ports - Get all port information
      POST /{moduleName}/api/v1/ports - Create new port
      GET /{moduleName}/api/v1/ports/{clientId}/{appName} - Get port information for specific client and application
      DELETE /{moduleName}/api/v1/ports?clientId=xx&appName=xx - Delete port
      Default port is 7890.
` + commonHelp

func main() {
	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	config := &chserver.Config{}
	flags.StringVar(&config.KeySeed, "key", "", "")
	flags.StringVar(&config.KeyFile, "keyfile", "", "")
	flags.StringVar(&config.AuthFile, "authfile", "", "")
	flags.StringVar(&config.Auth, "auth", "", "")
	flags.DurationVar(&config.KeepAlive, "keepalive", 25*time.Second, "")
	flags.StringVar(&config.Proxy, "proxy", "", "")
	flags.StringVar(&config.Proxy, "backend", "", "")
	flags.BoolVar(&config.Socks5, "socks5", false, "")
	flags.BoolVar(&config.Reverse, "reverse", false, "")
	flags.StringVar(&config.TLS.Key, "tls-key", "", "")
	flags.StringVar(&config.TLS.Cert, "tls-cert", "", "")
	flags.Var(multiFlag{&config.TLS.Domains}, "tls-domain", "")
	flags.StringVar(&config.TLS.CA, "tls-ca", "", "")
	flags.StringVar(&config.ControlPort, "control-port", "7890", "Control plane port, default is 7890")
	flags.IntVar(&config.MinPort, "minport", 30000, "Minimum port for port allocation, default is 30000")
	flags.IntVar(&config.MaxPort, "maxport", 31000, "Maximum port for port allocation, default is 31000")

	host := flags.String("host", "", "")
	p := flags.String("p", "", "")
	port := flags.String("port", "", "")
	pid := flags.Bool("pid", false, "")
	verbose := flags.Bool("v", false, "")
	keyGen := flags.String("keygen", "", "")

	flags.Usage = func() {
		fmt.Print(serverHelp)
		os.Exit(0)
	}
	flags.Parse(os.Args[1:])

	if *keyGen != "" {
		if err := ccrypto.GenerateKeyFile(*keyGen, config.KeySeed); err != nil {
			log.Fatal(err)
		}
		return
	}

	if config.KeySeed != "" {
		log.Print("Option `--key` is deprecated and will be removed in a future version of cotun.")
		log.Print("Please use `cotun server --keygen /file/path`, followed by `cotun server --keyfile /file/path` to specify the SSH private key")
	}

	if *host == "" {
		*host = os.Getenv("HOST")
	}
	if *host == "" {
		*host = "0.0.0.0"
	}
	if *port == "" {
		*port = *p
	}
	if *port == "" {
		*port = os.Getenv("PORT")
	}
	if *port == "" {
		*port = "8080"
	}
	if config.KeyFile == "" {
		config.KeyFile = settings.Env("KEY_FILE")
	} else if config.KeySeed == "" {
		config.KeySeed = settings.Env("KEY")
	}
	if config.Auth == "" {
		config.Auth = os.Getenv("AUTH")
	}
	s, err := chserver.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}
	s.Debug = *verbose
	if *pid {
		generatePidFile()
	}
	go cos.GoStats()
	ctx := cos.InterruptContext()
	if err := s.StartContext(ctx, *host, *port); err != nil {
		log.Fatal(err)
	}
	if err := s.Wait(); err != nil {
		log.Fatal(err)
	}
}

type multiFlag struct {
	values *[]string
}

func (flag multiFlag) String() string {
	return strings.Join(*flag.values, ", ")
}

func (flag multiFlag) Set(arg string) error {
	*flag.values = append(*flag.values, arg)
	return nil
}
