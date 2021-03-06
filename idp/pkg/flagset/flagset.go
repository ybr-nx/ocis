package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/idp/pkg/config"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Value:       "info",
			Usage:       "Set logging level",
			EnvVars:     []string{"IDP_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Value:       true,
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"IDP_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Value:       true,
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"IDP_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
	}
}

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9134",
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"IDP_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
	}
}

// ServerWithConfig applies cfg to the root flagset
func ServerWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "config-file",
			Value:       "",
			Usage:       "Path to config file",
			EnvVars:     []string{"IDP_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"IDP_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       "jaeger",
			Usage:       "Tracing backend type",
			EnvVars:     []string{"IDP_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       "",
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"IDP_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       "",
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"IDP_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       "idp",
			Usage:       "Service name for tracing",
			EnvVars:     []string{"IDP_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9134",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"IDP_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       "",
			Usage:       "Token to grant metrics access",
			EnvVars:     []string{"IDP_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"IDP_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVars:     []string{"IDP_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       "0.0.0.0:9130",
			Usage:       "Address to bind http server",
			EnvVars:     []string{"IDP_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "http-root",
			Value:       "/",
			Usage:       "Root path of http server",
			EnvVars:     []string{"IDP_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       "com.owncloud.web",
			Usage:       "Set the base namespace for service discovery",
			EnvVars:     []string{"IDP_HTTP_NAMESPACE"},
			Destination: &cfg.Service.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       "idp",
			Usage:       "Service name",
			EnvVars:     []string{"IDP_NAME"},
			Destination: &cfg.Service.Name,
		},
		&cli.StringFlag{
			Name:        "identity-manager",
			Value:       "ldap",
			Usage:       "Identity manager (one of ldap,kc,cookie,dummy)",
			EnvVars:     []string{"IDP_IDENTITY_MANAGER"},
			Destination: &cfg.IDP.IdentityManager,
		},
		&cli.StringFlag{
			Name:        "transport-tls-cert",
			Value:       "",
			Usage:       "Certificate file for transport encryption",
			EnvVars:     []string{"IDP_TRANSPORT_TLS_CERT"},
			Destination: &cfg.HTTP.TLSCert,
		},
		&cli.StringFlag{
			Name:        "transport-tls-key",
			Value:       "",
			Usage:       "Secret file for transport encryption",
			EnvVars:     []string{"IDP_TRANSPORT_TLS_KEY"},
			Destination: &cfg.HTTP.TLSKey,
		},
		&cli.StringFlag{
			Name:        "iss",
			Usage:       "OIDC issuer URL",
			EnvVars:     []string{"IDP_ISS", "OCIS_URL"}, // IDP_ISS takes precedence over OCIS_URL
			Value:       "https://localhost:9200",
			Destination: &cfg.IDP.Iss,
		},
		&cli.StringSliceFlag{
			Name:    "signing-private-key",
			Usage:   "Full path to PEM encoded private key file (must match the --signing-method algorithm)",
			EnvVars: []string{"IDP_SIGNING_PRIVATE_KEY"},
			Value:   nil,
		},
		&cli.StringFlag{
			Name:        "signing-kid",
			Usage:       "Value of kid field to use in created tokens (uniquely identifying the signing-private-key)",
			EnvVars:     []string{"IDP_SIGNING_KID"},
			Value:       "",
			Destination: &cfg.IDP.SigningKid,
		},
		&cli.StringFlag{
			Name:        "validation-keys-path",
			Usage:       "Full path to a folder containg PEM encoded private or public key files used for token validaton (file name without extension is used as kid)",
			EnvVars:     []string{"IDP_VALIDATION_KEYS_PATH"},
			Value:       "",
			Destination: &cfg.IDP.ValidationKeysPath,
		},
		&cli.StringFlag{
			Name:        "encryption-secret",
			Usage:       "Full path to a file containing a %d bytes secret key",
			EnvVars:     []string{"IDP_ENCRYPTION_SECRET"},
			Value:       "",
			Destination: &cfg.IDP.EncryptionSecretFile,
		},
		&cli.StringFlag{
			Name:        "signing-method",
			Usage:       "JWT default signing method",
			EnvVars:     []string{"IDP_SIGNING_METHOD"},
			Value:       "PS256",
			Destination: &cfg.IDP.SigningMethod,
		},
		&cli.StringFlag{
			Name:        "uri-base-path",
			Usage:       "Custom base path for URI endpoints",
			EnvVars:     []string{"IDP_URI_BASE_PATH"},
			Value:       "",
			Destination: &cfg.IDP.URIBasePath,
		},
		&cli.StringFlag{
			Name:        "sign-in-uri",
			Usage:       "Custom redirection URI to sign-in form",
			EnvVars:     []string{"IDP_SIGN_IN_URI"},
			Value:       "",
			Destination: &cfg.IDP.SignInURI,
		},
		&cli.StringFlag{
			Name:        "signed-out-uri",
			Usage:       "Custom redirection URI to signed-out goodbye page",
			EnvVars:     []string{"IDP_SIGN_OUT_URI"},
			Value:       "",
			Destination: &cfg.IDP.SignedOutURI,
		},
		&cli.StringFlag{
			Name:        "authorization-endpoint-uri",
			Usage:       "Custom authorization endpoint URI",
			EnvVars:     []string{"IDP_ENDPOINT_URI"},
			Value:       "",
			Destination: &cfg.IDP.AuthorizationEndpointURI,
		},
		&cli.StringFlag{
			Name:        "endsession-endpoint-uri",
			Usage:       "Custom endsession endpoint URI",
			EnvVars:     []string{"IDP_ENDSESSION_ENDPOINT_URI"},
			Value:       "",
			Destination: &cfg.IDP.EndsessionEndpointURI,
		},
		&cli.StringFlag{
			Name:        "asset-path",
			Value:       "",
			Usage:       "Path to custom assets",
			EnvVars:     []string{"IDP_ASSET_PATH"},
			Destination: &cfg.Asset.Path,
		},
		&cli.StringFlag{
			Name:        "identifier-client-path",
			Usage:       "Path to the identifier web client base folder",
			EnvVars:     []string{"IDP_IDENTIFIER_CLIENT_PATH"},
			Value:       "/var/tmp/ocis/idp",
			Destination: &cfg.IDP.IdentifierClientPath,
		},
		&cli.StringFlag{
			Name:        "identifier-registration-conf",
			Usage:       "Path to a identifier-registration.yaml configuration file",
			EnvVars:     []string{"IDP_IDENTIFIER_REGISTRATION_CONF"},
			Value:       "./config/identifier-registration.yaml",
			Destination: &cfg.IDP.IdentifierRegistrationConf,
		},
		&cli.StringFlag{
			Name:        "identifier-scopes-conf",
			Usage:       "Path to a scopes.yaml configuration file",
			EnvVars:     []string{"IDP_IDENTIFIER_SCOPES_CONF"},
			Value:       "",
			Destination: &cfg.IDP.IdentifierScopesConf,
		},
		&cli.BoolFlag{
			Name:        "insecure",
			Usage:       "Disable TLS certificate and hostname validation",
			EnvVars:     []string{"IDP_INSECURE"},
			Destination: &cfg.IDP.Insecure,
		},
		&cli.BoolFlag{
			Name:        "tls",
			Usage:       "Use TLS (disable only if idp is behind a TLS-terminating reverse-proxy).",
			EnvVars:     []string{"IDP_TLS"},
			Value:       false,
			Destination: &cfg.HTTP.TLS,
		},
		&cli.StringSliceFlag{
			Name:    "trusted-proxy",
			Usage:   "Trusted proxy IP or IP network (can be used multiple times)",
			EnvVars: []string{"IDP_TRUSTED_PROXY"},
			Value:   nil,
		},
		&cli.StringSliceFlag{
			Name:    "allow-scope",
			Usage:   "Allow OAuth 2 scope (can be used multiple times, if not set default scopes are allowed)",
			EnvVars: []string{"IDP_ALLOW_SCOPE"},
			Value:   nil,
		},
		&cli.BoolFlag{
			Name:        "allow-client-guests",
			Usage:       "Allow sign in of client controlled guest users",
			EnvVars:     []string{"IDP_ALLOW_CLIENT_GUESTS"},
			Destination: &cfg.IDP.AllowClientGuests,
		},
		&cli.BoolFlag{
			Name:        "allow-dynamic-client-registration",
			Usage:       "Allow dynamic OAuth2 client registration",
			EnvVars:     []string{"IDP_ALLOW_DYNAMIC_CLIENT_REGISTRATION"},
			Value:       true,
			Destination: &cfg.IDP.AllowDynamicClientRegistration,
		},
		&cli.BoolFlag{
			Name:        "disable-identifier-webapp",
			Usage:       "Disable built-in identifier-webapp to use a frontend hosted elsewhere.",
			EnvVars:     []string{"IDP_DISABLE_IDENTIFIER_WEBAPP"},
			Value:       true,
			Destination: &cfg.IDP.IdentifierClientDisabled,
		},
		&cli.Uint64Flag{
			Name:        "access-token-expiration",
			Usage:       "Expiration time of access tokens in seconds since generated",
			EnvVars:     []string{"IDP_ACCESS_TOKEN_EXPIRATION"},
			Destination: &cfg.IDP.AccessTokenDurationSeconds,
			Value:       60 * 10, // 10 Minutes.
		},
		&cli.Uint64Flag{
			Name:        "id-token-expiration",
			Usage:       "Expiration time of id tokens in seconds since generated",
			EnvVars:     []string{"IDP_ID_TOKEN_EXPIRATION"},
			Destination: &cfg.IDP.IDTokenDurationSeconds,
			Value:       60 * 60, // 1 Hour
		},
		&cli.Uint64Flag{
			Name:        "refresh-token-expiration",
			Usage:       "Expiration time of refresh tokens in seconds since generated",
			EnvVars:     []string{"IDP_REFRESH_TOKEN_EXPIRATION"},
			Destination: &cfg.IDP.RefreshTokenDurationSeconds,
			Value:       60 * 60 * 24 * 365 * 3, // 1 year
		},
	}
}

// ListIDPWithConfig applies the config to the list commands flags
func ListIDPWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{&cli.StringFlag{
		Name:        "http-namespace",
		Value:       "com.owncloud.web",
		Usage:       "Set the base namespace for service discovery",
		EnvVars:     []string{"IDP_HTTP_NAMESPACE"},
		Destination: &cfg.Service.Namespace,
	},
		&cli.StringFlag{
			Name:        "name",
			Value:       "idp",
			Usage:       "Service name",
			EnvVars:     []string{"IDP_NAME"},
			Destination: &cfg.Service.Name,
		},
	}
}
