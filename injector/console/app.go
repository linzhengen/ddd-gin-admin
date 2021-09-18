package console

type Options struct {
	ConfigFile string
	Version    string
}

type Option func(*Options)

func SetConfigFile(s string) Option {
	return func(o *Options) {
		o.ConfigFile = s
	}
}

func SetVersion(s string) Option {
	return func(o *Options) {
		o.Version = s
	}
}
