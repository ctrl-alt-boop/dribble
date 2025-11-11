package dsn

import (
	"net/url"
	"strings"

	"github.com/ctrl-alt-boop/dribble/datasource"
	"github.com/ctrl-alt-boop/dribble/internal/adapters/sql/sqlite3"
)

var _ datasource.Namer = (*SQLite3)(nil)

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

type SQLite3 struct {
	Path      string `json:"path"`
	AuthUser  string `json:"_auth_user"`
	AuthPass  string `json:"_auth_pass"`
	AuthCrypt string `json:"_auth_crypt"`
	Mode      Mode   `json:"mode"`
	Cache     Cache  `json:"cache"`
}

// SourceType implements datasource.Namer.
func (s *SQLite3) SourceType() datasource.SourceType {
	return sqlite3.SourceType
}

// Info implements database.DataSourceNamer.
func (s *SQLite3) Info() string {
	return "SQLite3: " + s.Path
}

// Type implements database.DataSourceNamer.
func (s SQLite3) Type() datasource.Type {
	return datasource.SQLite3
}

func (s SQLite3) DSN() string {
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

// SQLite3Option defines a function that configures a SQLite3 DSN.
type SQLite3Option func(*SQLite3)

// SQLite3DSN creates a new SQLite3DSN with the given options.
func SQLite3DSN(path string, opts ...SQLite3Option) *SQLite3 {
	// SQLite3 connection string format: file:test.db?cache=shared&mode=memory

	dsn := &SQLite3{
		Path: path,
	}
	for _, opt := range opts {
		opt(dsn)
	}
	return dsn
}

// SQLite3ReadOnly sets the DSN to read-only mode.
func SQLite3ReadOnly() SQLite3Option {
	return func(s *SQLite3) {
		s.Mode = ModeRO
	}
}

// SQLite3AuthUser sets the auth user for the DSN.
func SQLite3AuthUser(user string) SQLite3Option {
	return func(s *SQLite3) {
		s.AuthUser = user
	}
}

// SQLite3AuthPass sets the auth password for the DSN.
func SQLite3AuthPass(pass string) SQLite3Option {
	return func(s *SQLite3) {
		s.AuthPass = pass
	}
}

// SQLite3AuthCrypt sets the auth crypt for the DSN.
func SQLite3AuthCrypt(crypt string) SQLite3Option {
	return func(s *SQLite3) {
		s.AuthCrypt = crypt
	}
}

// SQLite3Mode sets the mode for the DSN.
func SQLite3Mode(mode Mode) SQLite3Option {
	return func(s *SQLite3) {
		s.Mode = mode
	}
}

// SQLite3Cache sets the cache for the DSN.
func SQLite3Cache(cache Cache) SQLite3Option {
	return func(s *SQLite3) {
		s.Cache = cache
	}
}
