package featureflag

import "testing"

func TestEnabledFeature(t *testing.T) {
	name := "TestEnabledFeature"
	enabled(name)
	if !features[name].Enabled() {
		t.Errorf("got %s.Enabled() = false, want true.", name)
	}
}

func TestDisabledFeature(t *testing.T) {
	name := "TestDisabledFeature"
	disabled(name)
	if features[name].Enabled() {
		t.Errorf("got %s.Enabled() = true, want false.", name)
	}
}

func TestEnabledFeatureEnv(t *testing.T) {
	name := "TestEnabledFeatureEnv"
	disabled(name)
	t.Setenv("DEVBOX_FEATURE_"+name, "1")
	if !features[name].Enabled() {
		t.Errorf("got %s.Enabled() = false, want true.", name)
	}
}

func TestNonExistentFeature(t *testing.T) {
	name := "TestNonExistentFeature"
	if features[name].Enabled() {
		t.Errorf("got %s.Enabled() = true, want false.", name)
	}
}
