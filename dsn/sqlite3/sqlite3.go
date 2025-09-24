package sqlite3

import (
	"net/url"
	"strings"

	"github.com/ctrl-alt-boop/dribble/database"
)

var _ database.DataSourceNamer = (*SQLite3DSN)(nil)

type Mode string

const (
	ModeRO  Mode = "ro"
	ModeRW  Mode = "rw"
	ModeRWC Mode = "rwc"
	ModeMem Mode = "memory"
)

type Cache string

const (
	CacheShared  Cache = "shared"
	CachePrivate Cache = "private"
)

type SQLite3DSN struct {
	Path      string `json:"path"`
	AuthUser  string `json:"_auth_user"`
	AuthPass  string `json:"_auth_pass"`
	AuthCrypt string `json:"_auth_crypt"`
	Mode      Mode   `json:"mode"`
	Cache     Cache  `json:"cache"`
}

// Info implements database.DataSourceNamer.
func (s *SQLite3DSN) Info() string {
	return "SQLite3: " + s.Path
}

// Type implements database.DataSourceNamer.
func (s SQLite3DSN) Type() database.Type {
	return database.SQLite3
}

func (s SQLite3DSN) DSN() string {
	if s.Path == "" {
		return ""
	}

	var queryParts []string

	if s.Mode != "" {
		queryParts = append(queryParts, "mode="+string(s.Mode))
	}
	if s.Cache != "" {
		queryParts = append(queryParts, "cache="+string(s.Cache))
	}

	hasAuth := s.AuthUser != "" || s.AuthPass != "" || s.AuthCrypt != ""
	if hasAuth {
		queryParts = append(queryParts, "_auth")
		if s.AuthUser != "" {
			queryParts = append(queryParts, "_auth_user="+url.QueryEscape(s.AuthUser))
		}
		if s.AuthPass != "" {
			queryParts = append(queryParts, "_auth_pass="+url.QueryEscape(s.AuthPass))
		}
		if s.AuthCrypt != "" {
			queryParts = append(queryParts, "_auth_crypt="+url.QueryEscape(s.AuthCrypt))
		}
	}

	if len(queryParts) == 0 {
		return "file:" + s.Path
	}

	return "file:" + s.Path + "?" + strings.Join(queryParts, "&")
}

// DSNOption defines a function that configures a SQLite3DSN.
type DSNOption func(*SQLite3DSN)

// NewDSN creates a new SQLite3DSN with the given options.
func NewDSN(path string, opts ...DSNOption) *SQLite3DSN {
	// SQLite3 connection string format: file:test.db?cache=shared&mode=memory

	dsn := &SQLite3DSN{
		Path: path,
	}
	for _, opt := range opts {
		opt(dsn)
	}
	return dsn
}

func ReadOnly() DSNOption {
	return func(s *SQLite3DSN) {
		s.Mode = ModeRO
	}
}

// WithAuthUser sets the auth user for the DSN.
func WithAuthUser(user string) DSNOption {
	return func(s *SQLite3DSN) {
		s.AuthUser = user
	}
}

// WithAuthPass sets the auth password for the DSN.
func WithAuthPass(pass string) DSNOption {
	return func(s *SQLite3DSN) {
		s.AuthPass = pass
	}
}

// WithAuthCrypt sets the auth crypt for the DSN.
func WithAuthCrypt(crypt string) DSNOption {
	return func(s *SQLite3DSN) {
		s.AuthCrypt = crypt
	}
}

// WithMode sets the mode for the DSN.
func WithMode(mode Mode) DSNOption {
	return func(s *SQLite3DSN) {
		s.Mode = mode
	}
}

// WithCache sets the cache for the DSN.
func WithCache(cache Cache) DSNOption {
	return func(s *SQLite3DSN) {
		s.Cache = cache
	}
}
