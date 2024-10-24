package bindings

import (
	"net/url"

	"code.cloudfoundry.org/loggregator-agent-release/src/pkg/binding"
	"code.cloudfoundry.org/loggregator-agent-release/src/pkg/egress/syslog"
)

type DrainParamParser struct {
	fetcher              binding.Fetcher
	defaultDrainMetadata bool
}

func NewDrainParamParser(f binding.Fetcher, defaultDrainMetadata bool) *DrainParamParser {
	return &DrainParamParser{
		fetcher:              f,
		defaultDrainMetadata: defaultDrainMetadata,
	}
}

func (d *DrainParamParser) FetchBindings() ([]syslog.Binding, error) {
	var processed []syslog.Binding
	bs, err := d.fetcher.FetchBindings()
	if err != nil {
		return nil, err
	}

	for _, b := range bs {
		urlParsed, err := url.Parse(b.Drain.Url)
		if err != nil {
			continue
		}

		if d.defaultDrainMetadata && urlParsed.Query().Get("disable-metadata") == "true" {
			b.OmitMetadata = true
		} else if !d.defaultDrainMetadata && urlParsed.Query().Get("disable-metadata") == "false" {
			b.OmitMetadata = false
		} else {
			b.OmitMetadata = !d.defaultDrainMetadata
		}
		if urlParsed.Query().Get("ssl-strict-internal") == "true" {
			b.InternalTls = true
		}

		processed = append(processed, b)
	}

	return processed, nil
}

func (d *DrainParamParser) DrainLimit() int {
	return -1
}
