/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fake

import (
	"fmt"

	"openshift/node-feature-discovery/source"
)

// Configuration file options
type Config struct {
	Labels map[string]string `json:"labels"`
}

// newDefaultConfig returns a new config with defaults values
func newDefaultConfig() *Config {
	return &Config{
		Labels: map[string]string{
			"fakefeature1": "true",
			"fakefeature2": "true",
			"fakefeature3": "true",
		},
	}
}

// Source implements FeatureSource.
type Source struct {
	config *Config
}

// Name returns an identifier string for this feature source.
func (s Source) Name() string { return "fake" }

// NewConfig method of the FeatureSource interface
func (s *Source) NewConfig() source.Config { return newDefaultConfig() }

// GetConfig method of the FeatureSource interface
func (s *Source) GetConfig() source.Config { return s.config }

// SetConfig method of the FeatureSource interface
func (s *Source) SetConfig(conf source.Config) {
	switch v := conf.(type) {
	case *Config:
		s.config = v
	default:
		panic(fmt.Sprintf("invalid config type: %T", conf))
	}
}

// Configure method of the FeatureSource interface
func (s Source) Configure([]byte) error { return nil }

// Discover returns feature names for some fake features.
func (s Source) Discover() (source.Features, error) {
	// Adding three fake features.
	features := make(source.Features, len(s.config.Labels))
	for k, v := range s.config.Labels {
		features[k] = v
	}

	return features, nil
}
