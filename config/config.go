package config

import (
	"os"

	"github.com/jinzhu/configor"
	"github.com/qor/render"
)

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Site     string
}

var Config = struct {
	Port uint   `default:"80" env:"PORT"`
	Host string `default:"192.168.1.5" env:"HOST"`
	DB   struct {
		Name     string `default:"qor_example"`
		Adapter  string `default:"mysql"`
		User     string
		Password string
	}
	SMTP SMTPConfig
}{}

var (
	Root = os.Getenv("GOPATH") + "/src/github.com/sunwukonga/qor-example"
	View *render.Render
)

func init() {
	if err := configor.Load(&Config, "config/database.yml", "config/smtp.yml"); err != nil {
		panic(err)
	}

	View = render.New()
}

func (s SMTPConfig) HostWithPort() string {
	return s.Host + ":" + s.Port
}
