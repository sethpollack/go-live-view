package uploads

import "mime"

type Option func(*Config)

func WithExternal(f func() error) Option {
	return func(u *Config) {
		u.External = true
		u.PresignFunc = f
	}
}

func WithAccept(accept ...string) Option {
	return func(u *Config) {
		for _, a := range accept {
			u.Accept = append(u.Accept, a)
			u.MimeTypes = append(u.MimeTypes, mime.TypeByExtension(a))
		}
	}
}

func WithAutoUpload(auto bool) Option {
	return func(u *Config) {
		u.AutoUpload = auto
	}
}

func WithMaxEntries(max int) Option {
	return func(u *Config) {
		u.MaxEntries = max
	}
}

func WithMaxFileSize(max int) Option {
	return func(u *Config) {
		u.MaxFileSize = max
	}
}

func WithChunkSize(size int) Option {
	return func(u *Config) {
		u.ChunkSize = size
	}
}

func WithChunkTimeout(timeout int) Option {
	return func(u *Config) {
		u.ChunkTimeout = timeout
	}
}

func WithWriter(w Writer) Option {
	return func(u *Config) {
		u.Writer = w
	}
}
