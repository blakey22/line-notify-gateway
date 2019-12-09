package flag

var Options struct {
	Host      string `short:"h" long:"host" description:"Bind host" env:"HOST" default:"0.0.0.0"`
	Port      int    `short:"p" long:"port" description:"Bind port" env:"PORT" default:"8080"`
	Secret    string `short:"s" long:"secret" description:"Secret token for accessing the gateway" env:"SECRET" default:""`
	Token     string `short:"t" long:"token" description:"LINE notify token" env:"TOKEN" default:""`
	Endpoint  string `long:"endpoint" description:"LINE notify endpoint" env:"ENDPOINT" default:"https://notify-api.line.me/api/notify"`
	Locale    string `long:"locale" description:"Locale of LINE message" env:"LOCALE" default:""`
	Templates string `long:"templates-path" description:"Path of message templates" env:"TEMPLATES_PATH" default:"./templates"`
}
