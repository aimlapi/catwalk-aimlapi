package providers

import (
	"regexp"
	"slices"
	"strings"
	"testing"

	"charm.land/catwalk/pkg/catwalk"
)

func TestValidDefaultModels(t *testing.T) {
	for _, p := range GetAll() {
		t.Run(p.Name, func(t *testing.T) {
			var modelIds []string
			for _, m := range p.Models {
				modelIds = append(modelIds, m.ID)
			}
			if !slices.Contains(modelIds, p.DefaultLargeModelID) {
				t.Errorf("Default large model %q not found in provider %q", p.DefaultLargeModelID, p.Name)
			}
			if !slices.Contains(modelIds, p.DefaultSmallModelID) {
				t.Errorf("Default small model %q not found in provider %q", p.DefaultSmallModelID, p.Name)
			}
		})
	}
}

// aimlapiPartnerIDPattern is the shape the aimlapi.com gateway accepts for
// X-AIMLAPI-Partner-ID. A value that does not match is silently ignored
// server-side, so a typo here would never surface at runtime.
var aimlapiPartnerIDPattern = regexp.MustCompile(`^part_[A-Za-z0-9]{1,64}$`)

func TestAIMLAPIDefaultHeaders(t *testing.T) {
	var found bool
	for _, p := range GetAll() {
		if p.ID != catwalk.InferenceProviderAIMLAPI {
			// Attribution headers must never ride along to another provider.
			for k := range p.DefaultHeaders {
				if strings.HasPrefix(strings.ToLower(k), "x-aimlapi-") {
					t.Errorf("provider %q must not send %q", p.ID, k)
				}
			}
			continue
		}
		found = true
		if got := p.DefaultHeaders["X-AIMLAPI-Partner-ID"]; !aimlapiPartnerIDPattern.MatchString(got) {
			t.Errorf("X-AIMLAPI-Partner-ID %q does not match %s", got, aimlapiPartnerIDPattern)
		}
		if got, want := p.DefaultHeaders["X-AIMLAPI-Source"], "agent/crush"; got != want {
			t.Errorf("X-AIMLAPI-Source = %q, want %q", got, want)
		}
	}
	if !found {
		t.Fatalf("provider %q not registered", catwalk.InferenceProviderAIMLAPI)
	}
}
