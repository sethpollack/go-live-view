package uploads

import (
	"fmt"
	"go-live-view/internal/ref"
	"go-live-view/params"
	"strings"
	"sync"
)

type Uploads struct {
	mu      sync.RWMutex
	ref     *ref.Ref
	refMap  map[string]string
	uploads map[string]*Config
}

type Config struct {
	Ref          string
	Name         string
	Accept       []string
	MimeTypes    []string
	AutoUpload   bool
	MaxEntries   int
	MaxFileSize  int
	ChunkSize    int
	ChunkTimeout int
	Entries      []*Entry
	Errors       []string
	Writer       Writer
	External     bool
	PresignFunc  func() error
}

type Meta struct {
	RelativePath string
	Size         int
	LastModified int
	FileType     string
}

type Entry struct {
	ConfigRef string
	Ref       string

	UUID string

	Meta Meta

	Errors   []string
	Progress float32

	Valid     bool
	Preflight bool
	Cancelled bool
	Done      bool

	closeClient func() error
}

func New() *Uploads {
	return &Uploads{
		ref:     ref.New(0),
		refMap:  make(map[string]string),
		uploads: make(map[string]*Config),
	}
}

func (u *Uploads) GetByName(name string) *Config {
	u.mu.RLock()
	defer u.mu.RUnlock()

	ref, ok := u.refMap[name]
	if ok {
		return u.GetByRef(ref)
	}

	return nil
}

func (u *Uploads) GetByRef(ref string) *Config {
	u.mu.RLock()
	defer u.mu.RUnlock()

	upload, ok := u.uploads[ref]
	if !ok {
		return nil
	}

	return upload
}

func (u *Uploads) AllowUpload(name string, options ...Option) {
	u.mu.Lock()
	defer u.mu.Unlock()

	ref := u.ref.NextStringRef()

	u.refMap[name] = ref

	c := &Config{
		Name:         name,
		Ref:          ref,
		MaxEntries:   1,
		MaxFileSize:  8 * 1024 * 1024,
		ChunkSize:    64 * 1024,
		ChunkTimeout: 10 * 1000,
		Accept:       []string{"*"},
		Writer:       &TmpFileWriter{},
	}

	for _, option := range options {
		option(c)
	}

	u.uploads[ref] = c
}

func (u *Uploads) Consume(name string, f func(path string)) error {
	cfg := u.GetByName(name)
	if cfg == nil {
		return fmt.Errorf("upload not found")
	}
	for _, entry := range cfg.Entries {
		if entry.Done {
			f(entry.Meta.RelativePath)
		}
	}

	cfg.reset()

	return nil
}

func (u *Uploads) Cancel(name string, ref string) error {
	cfg := u.GetByName(name)
	if cfg == nil {
		return fmt.Errorf("upload not found")
	}

	for _, entry := range cfg.Entries {
		if entry.Ref == ref {
			entry.Cancelled = true
			if entry.closeClient != nil {
				entry.closeClient()
			}
		}
	}
	return nil
}

func (u *Uploads) OnValidate(params params.Params) {
	upload := params.Map("uploads")
	for ref := range upload {
		c, ok := u.uploads[ref]
		if ok {
			entries := upload.Slice(ref)
			for _, entry := range entries {
				c.Entries = append(c.Entries, &Entry{
					Ref: entry.String("ref"),
					Meta: Meta{
						FileType:     entry.String("file_type"),
						LastModified: entry.Int("last_modified"),
						RelativePath: entry.String("relative_path"),
						Size:         entry.Int("size"),
					},
				})
			}
			c.validate()
		}
	}
}

func (c *Config) OnAllowUploads(params params.Params) {
	entries := params.Slice("entries")

	c.Entries = make([]*Entry, len(entries))
	for i, entry := range entries {
		c.Entries[i] = &Entry{
			Ref: entry.String("ref"),
			Meta: Meta{
				FileType:     entry.String("file_type"),
				LastModified: entry.Int("last_modified"),
				RelativePath: entry.String("relative_path"),
				Size:         entry.Int("size"),
			},
			UUID:      c.Ref + "-" + entry.String("ref"), // TODO: encode with proper token
			Preflight: true,
		}
	}

	c.validate()
}

func (c *Config) OnChunk(ref string, data []byte, close func() error) error {
	for _, entry := range c.Entries {
		if entry.Ref != ref {
			continue
		}

		entry.closeClient = close

		_, err := c.Writer.WriteChunk(data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) OnProgress(ref string, progress float32) error {
	for _, entry := range c.Entries {
		if entry.Ref == ref {
			entry.Progress = progress
			if progress == 100 {
				entry.Done = true
				if entry.closeClient != nil {
					err := entry.closeClient()
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (c *Config) PreflightEntries() map[string]string {
	entries := make(map[string]string)
	for _, entry := range c.Entries {
		if entry.Preflight {
			entries[entry.Ref] = entry.UUID
		}
	}

	return entries
}

func (c *Config) PreflightErrors() [][]string {
	errors := [][]string{}

	for _, err := range c.Errors {
		errors = append(errors, []string{c.Ref, err})
	}

	for _, entry := range c.Entries {
		for _, err := range entry.Errors {
			errors = append(errors, []string{entry.Ref, err})
		}
	}

	return errors
}

func (c *Config) PreflightRefs() *string {
	return getRefs(c, func(e *Entry) bool {
		return e.Preflight
	})
}

func (c *Config) ActiveRefs() *string {
	return getRefs(c, func(e *Entry) bool {
		return true
	})
}

func (c *Config) DoneRefs() *string {
	return getRefs(c, func(e *Entry) bool {
		return e.Done
	})
}

func (c *Config) validate() {
	if len(c.Entries) > c.MaxEntries {
		c.Errors = append(c.Errors, "Max entries exceeded")
	}

	for _, entry := range c.Entries {
		if entry == nil {
			continue
		}
		if c.MaxFileSize > 0 && entry.Meta.Size > c.MaxFileSize {
			entry.Errors = append(entry.Errors, "Max file size exceeded")
			entry.Valid = false
		}

		if contains(c.Accept, "*") {
			continue
		}

		if !contains(c.MimeTypes, entry.Meta.FileType) {
			entry.Errors = append(entry.Errors, "Invalid file type")
			entry.Valid = false
		}
	}
}

func (c *Config) reset() {
	c.Entries = nil
}

func getRefs(c *Config, f func(e *Entry) bool) *string {
	refs := []string{}
	for _, entry := range c.Entries {
		if f(entry) {
			refs = append(refs, entry.Ref)
		}
	}

	res := strings.Join(refs, ",")

	return &res
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
