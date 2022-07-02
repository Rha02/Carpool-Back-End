package config

import "github.com/gorilla/sessions"

type AppConfig struct {
	CookieStore *sessions.CookieStore
}
