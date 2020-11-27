package cs3users

import (
	idxerrs "github.com/owncloud/ocis/ocis-pkg/indexer/errors"
	"github.com/owncloud/ocis/ocis-pkg/indexer/index"
	"github.com/owncloud/ocis/ocis-pkg/indexer/option"
	"github.com/owncloud/ocis/ocis-pkg/indexer/registry"
)

// Autoincrement are fields for an index of type autoincrement.
type Autoincrement struct {
	indexBy  string
	typeName string

	bound *option.Bound
}

func init() {
	registry.IndexConstructorRegistry["cs3users"]["autoincrement"] = NewAutoincrementIndex
}

// NewAutoincrementIndex instantiates a new AutoincrementIndex instance. Init() MUST be called upon instantiation.
func NewAutoincrementIndex(o ...option.Option) index.Index {
	opts := &option.Options{}
	for _, opt := range o {
		opt(opts)
	}

	if opts.Entity == nil {
		panic("invalid autoincrement index: configured without entity")
	}

	k, err := getKind(opts.Entity, opts.IndexBy)
	if !isValidKind(k) || err != nil {
		panic("invalid autoincrement index: configured on non-numeric field")
	}

	return &Autoincrement{
		indexBy:  opts.IndexBy,
		typeName: opts.TypeName,
		bound:    opts.Bound,
	}
}

// Init initializes an autoincrement index.
func (idx *Autoincrement) Init() error {
	return nil
}

// Lookup exact lookup by value.
func (idx *Autoincrement) Lookup(v string) ([]string, error) {
	return []string{}, nil
}

// Add a new value to the index.
func (idx *Autoincrement) Add(id, v string) (string, error) {
	return "", nil
}

// Remove a value v from an index.
func (idx *Autoincrement) Remove(id string, v string) error {
	return &idxerrs.NotSupportedErr{}
}

// Update index from <oldV> to <newV>.
func (idx *Autoincrement) Update(id, oldV, newV string) error {
	return &idxerrs.NotSupportedErr{}
}

// Search allows for glob search on the index.
func (idx *Autoincrement) Search(pattern string) ([]string, error) {
	// TODO implement search by uid
	return []string{}, nil
}

// CaseInsensitive undocumented.
func (idx *Autoincrement) CaseInsensitive() bool {
	return false
}

// IndexBy undocumented.
func (idx *Autoincrement) IndexBy() string {
	return idx.indexBy
}

// TypeName undocumented.
func (idx *Autoincrement) TypeName() string {
	return idx.typeName
}

// FilesDir  undocumented.
func (idx *Autoincrement) FilesDir() string {
	return ""
}

// Delete deletes the index root folder from the configured storage.
func (idx *Autoincrement) Delete() error {
	return &idxerrs.NotSupportedErr{}
}
