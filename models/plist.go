package models

import (
	"io"
	"strings"

	"github.com/DHowett/go-plist"
)

const (
	PlistFileName                   = "test.plist"
	AssetKind                       = "software-package"
	DefaultMetadataBundleIdentifier = "com.example.test"
	MetadataKind                    = "software"
)

type Plist struct {
	Items []*Item `plist:"items"`
}

type Item struct {
	Assets   []*Asset  `plist:"assets"`
	Metadata *Metadata `plist:"metadata"`
}

type Asset struct {
	Kind string `plist:"kind"`
	Url  string `plist:"url"`
}

type Metadata struct {
	BundleIdentifier string `plist:"bundle-identifier"`
	BundleVersion    string `plist:"bundle-version"`
	Kind             string `plist:"kind"`
	Title            string `plist:"title"`
}

func NewPlist(title, version, identifier, ipaUrl string) *Plist {
	if len(identifier) == 0 {
		identifier = DefaultMetadataBundleIdentifier
	}

	return &Plist{
		Items: []*Item{
			&Item{
				Assets: []*Asset{
					&Asset{
						Kind: AssetKind,
						Url:  ipaUrl,
					},
				},
				Metadata: &Metadata{
					BundleIdentifier: identifier,
					BundleVersion:    version,
					Kind:             MetadataKind,
					Title:            title,
				},
			},
		},
	}
}

func (p *Plist) Marshall() ([]byte, error) {
	return plist.MarshalIndent(p, plist.XMLFormat, "\t")
}

func (p *Plist) Reader() (io.Reader, error) {
	data, err := p.Marshall()
	if err != nil {
		return nil, err
	}

	return strings.NewReader(string(data)), nil
}
