package svc

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-chi/chi"
	"github.com/gorilla/mux"
	"github.com/owncloud/ocis/konnectd/pkg/assets"
	"github.com/owncloud/ocis/konnectd/pkg/config"
	logw "github.com/owncloud/ocis/konnectd/pkg/log"
	"github.com/owncloud/ocis/konnectd/pkg/middleware"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"stash.kopano.io/kc/konnect/bootstrap"
	kcconfig "stash.kopano.io/kc/konnect/config"
	"stash.kopano.io/kc/konnect/server"
	"stash.kopano.io/kgol/rndm"
)

// Service defines the extension handlers.
type Service interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
	Dummy(http.ResponseWriter, *http.Request)
}

// NewService returns a service implementation for Service.
func NewService(opts ...Option) Service {
	ctx := context.Background()
	options := newOptions(opts...)
	logger := options.Logger.Logger
	assetVFS := assets.New(
		assets.Logger(options.Logger),
		assets.Config(options.Config),
	)

	if err := initKonnectInternalEnvVars(); err != nil {
		logger.Fatal().Err(err).Msg("could not initialize env vars")
	}

	if err := createConfigsIfNotExist(options.Config, assetVFS); err != nil {
		logger.Fatal().Err(err).Msg("could not create default config")
	}

	bs, err := bootstrap.Boot(ctx, &options.Config.Konnectd, &kcconfig.Config{
		Logger: logw.Wrap(logger),
	})

	if err != nil {
		logger.Fatal().Err(err).Msg("could not bootstrap konnectd")
	}

	managers := bs.Managers()
	routes := []server.WithRoutes{managers.Must("identity").(server.WithRoutes)}
	handlers := managers.Must("handler").(http.Handler)

	svc := Konnectd{
		logger: options.Logger,
		config: options.Config,
		assets: assetVFS,
	}

	svc.initMux(ctx, routes, handlers, options)

	return svc
}

func createConfigsIfNotExist(cfg *config.Config, assets http.FileSystem) error {

	configPath := filepath.Dir(cfg.Konnectd.IdentifierRegistrationConf)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.Mkdir(configPath, 0700); err != nil {
			return err
		}
	}

	defaultConf, err := assets.Open("/identifier-registration.yaml")
	if err != nil {
		return err
	}
	defer defaultConf.Close()

	if _, err := os.Stat(cfg.Konnectd.IdentifierRegistrationConf); os.IsNotExist(err) {
		confOnDisk, err := os.Create(cfg.Konnectd.IdentifierRegistrationConf)
		if err != nil {
			return err
		}

		defer confOnDisk.Close()

		conf, err := ioutil.ReadAll(defaultConf)
		if err != nil {
			return err
		}

		conf, _ = templateConfig(cfg, conf)

		err = ioutil.WriteFile(cfg.Konnectd.IdentifierRegistrationConf, conf, 0600)
		if err != nil {
			return err
		}
	} else {
		confOnDisk, err := os.Open(cfg.Konnectd.IdentifierRegistrationConf)
		if err != nil {
			return err
		}

		defer confOnDisk.Close()

		conf, err := ioutil.ReadAll(confOnDisk)
		if err != nil {
			return err
		}

		re := regexp.MustCompile(`sha256: (.*?)\n`)
		matches := re.FindStringSubmatch(string(conf))

		checksum := ""
		if len(matches) > 1 {
			checksum = matches[1]
		}

		if checksum != confChecksum(conf) {
			// The config was touched meanwhile. Do not change it!
			return nil
		}

		conf, err = ioutil.ReadAll(defaultConf)
		if err != nil {
			return err
		}

		conf, _ = templateConfig(cfg, conf)

		err = ioutil.WriteFile(cfg.Konnectd.IdentifierRegistrationConf, conf, 0600)
		if err != nil {
			return err
		}

	}

	return nil

}

func templateConfig(cfg *config.Config, conf []byte) (config []byte, checksum string) {
	// replace placeholder {{OCIS_URL}} with https://localhost:9200 / correct host
	conf = []byte(strings.ReplaceAll(string(conf), "{{OCIS_URL}}", strings.TrimRight(cfg.Konnectd.Iss, "/")))

	checksum = confChecksum(conf)

	checksumConf := append([]byte("sha256: "), checksum[:]...)
	checksumConf = append(checksumConf, []byte("\n")...)
	checksumConf = append(checksumConf, conf...)

	return checksumConf, checksum
}

func confChecksum(conf []byte) string {
	if !strings.Contains(string(conf), "---\n") {
		conf = append([]byte("---\n"), conf...)
	}

	sha := sha256.Sum256([]byte(strings.Split(string(conf), "---")[1]))
	checksum := hex.EncodeToString(sha[:])
	return checksum
}

// Init vars which are currently not accessible via konnectd api
func initKonnectInternalEnvVars() error {
	var defaults = map[string]string{
		"LDAP_URI":                 "ldap://localhost:9125",
		"LDAP_BINDDN":              "cn=konnectd,ou=sysusers,dc=example,dc=org",
		"LDAP_BINDPW":              "konnectd",
		"LDAP_BASEDN":              "ou=users,dc=example,dc=org",
		"LDAP_SCOPE":               "sub",
		"LDAP_LOGIN_ATTRIBUTE":     "cn",
		"LDAP_EMAIL_ATTRIBUTE":     "mail",
		"LDAP_NAME_ATTRIBUTE":      "sn",
		"LDAP_UUID_ATTRIBUTE":      "uid",
		"LDAP_UUID_ATTRIBUTE_TYPE": "text",
		"LDAP_FILTER":              "(objectClass=posixaccount)",
	}

	for k, v := range defaults {
		if _, exists := os.LookupEnv(k); !exists {
			if err := os.Setenv(k, v); err != nil {
				return fmt.Errorf("could not set env var %s=%s", k, v)
			}
		}
	}

	return nil
}

// Konnectd defines implements the business logic for Service.
type Konnectd struct {
	logger log.Logger
	config *config.Config
	mux    *chi.Mux
	assets http.FileSystem
}

// initMux initializes the internal konnectd gorilla mux and mounts it in to a ocis chi-router
func (k *Konnectd) initMux(ctx context.Context, r []server.WithRoutes, h http.Handler, options Options) {
	gm := mux.NewRouter()
	for _, route := range r {
		route.AddRoutes(ctx, gm)
	}

	// Delegate rest to provider which is also a handler.
	if h != nil {
		gm.NotFoundHandler = h
	}

	k.mux = chi.NewMux()
	k.mux.Use(options.Middleware...)

	k.mux.Use(middleware.Static(
		"/signin/v1/",
		assets.New(
			assets.Logger(options.Logger),
			assets.Config(options.Config),
		),
	))

	// handle / | index.html with a template that needs to have the BASE_PREFIX replaced
	k.mux.Get("/signin/v1/identifier", k.Index())
	k.mux.Get("/signin/v1/identifier/", k.Index())
	k.mux.Get("/signin/v1/identifier/index.html", k.Index())

	k.mux.Mount("/", gm)
}

// ServeHTTP implements the Service interface.
func (k Konnectd) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	k.mux.ServeHTTP(w, r)
}

// Dummy implements the Service interface.
func (k Konnectd) Dummy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// Index renders the static html with the
func (k Konnectd) Index() http.HandlerFunc {

	f, err := k.assets.Open("/identifier/index.html")
	if err != nil {
		k.logger.Fatal().Err(err).Msg("Could not open index template")
	}

	template, err := ioutil.ReadAll(f)
	if err != nil {
		k.logger.Fatal().Err(err).Msg("Could not read index template")
	}
	if err = f.Close(); err != nil {
		k.logger.Fatal().Err(err).Msg("Could not close body")
	}

	// TODO add environment variable to make the path prefix configurable
	pp := "/signin/v1"
	indexHTML := bytes.Replace(template, []byte("__PATH_PREFIX__"), []byte(pp), 1)

	nonce := rndm.GenerateRandomString(32)
	indexHTML = bytes.Replace(indexHTML, []byte("__CSP_NONCE__"), []byte(nonce), 1)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(indexHTML)
	})
}
