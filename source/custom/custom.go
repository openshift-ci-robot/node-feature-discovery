/*
Copyright 2020-2021 The Kubernetes Authors.

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

package custom

import (
	"reflect"

	"k8s.io/klog/v2"

	"openshift/node-feature-discovery/pkg/utils"
	"openshift/node-feature-discovery/source"
	"openshift/node-feature-discovery/source/custom/rules"
)

// Custom Features Configurations
type MatchRule struct {
	PciID      *rules.PciIDRule      `json:"pciId,omitempty"`
	UsbID      *rules.UsbIDRule      `json:"usbId,omitempty"`
	LoadedKMod *rules.LoadedKModRule `json:"loadedKMod,omitempty"`
	CpuID      *rules.CpuIDRule      `json:"cpuId,omitempty"`
	Kconfig    *rules.KconfigRule    `json:"kConfig,omitempty"`
	Nodename   *rules.NodenameRule   `json:"nodename,omitempty"`
}

type FeatureSpec struct {
	Name    string      `json:"name"`
	Value   *string     `json:"value,omitempty"`
	MatchOn []MatchRule `json:"matchOn"`
}

type config []FeatureSpec

// newDefaultConfig returns a new config with pre-populated defaults
func newDefaultConfig() *config {
	return &config{}
}

// Implements FeatureSource Interface
type Source struct {
	config *config
}

// Return name of the feature source
func (s Source) Name() string { return "custom" }

// NewConfig method of the FeatureSource interface
func (s *Source) NewConfig() source.Config { return newDefaultConfig() }

// GetConfig method of the FeatureSource interface
func (s *Source) GetConfig() source.Config { return s.config }

// SetConfig method of the FeatureSource interface
func (s *Source) SetConfig(conf source.Config) {
	switch v := conf.(type) {
	case *config:
		s.config = v
	default:
		klog.Fatalf("invalid config type: %T", conf)
	}
}

// Discover features
func (s Source) Discover() (source.Features, error) {
	features := source.Features{}
	allFeatureConfig := append(getStaticFeatureConfig(), *s.config...)
	allFeatureConfig = append(allFeatureConfig, getDirectoryFeatureConfig()...)
	utils.KlogDump(2, "custom features configuration:", "  ", allFeatureConfig)
	// Iterate over features
	for _, customFeature := range allFeatureConfig {
		featureExist, err := s.discoverFeature(customFeature)
		if err != nil {
			klog.Errorf("failed to discover feature: %q: %s", customFeature.Name, err.Error())
			continue
		}
		if featureExist {
			var value interface{} = true
			if customFeature.Value != nil {
				value = *customFeature.Value
			}
			features[customFeature.Name] = value
		}
	}
	return features, nil
}

// Process a single feature by Matching on the defined rules.
// A feature is present if all defined Rules in a MatchRule return a match.
func (s Source) discoverFeature(feature FeatureSpec) (bool, error) {
	for _, matchRules := range feature.MatchOn {

		allRules := []rules.Rule{
			matchRules.PciID,
			matchRules.UsbID,
			matchRules.LoadedKMod,
			matchRules.CpuID,
			matchRules.Kconfig,
			matchRules.Nodename,
		}

		// return true, nil if all rules match
		matchRules := func(rules []rules.Rule) (bool, error) {
			for _, rule := range rules {
				if reflect.ValueOf(rule).IsNil() {
					continue
				}
				if match, err := rule.Match(); err != nil {
					return false, err
				} else if !match {
					return false, nil
				}
			}
			return true, nil
		}

		if match, err := matchRules(allRules); err != nil {
			return false, err
		} else if match {
			return true, nil
		}
	}
	return false, nil
}
