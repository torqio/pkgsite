// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestValidateAppVersion(t *testing.T) {
	for _, tc := range []struct {
		in      string
		wantErr bool
	}{
		{"", true},
		{"20190912t130708", false},
		{"20190912t130708x", true},
		{"2019-09-12t13-07-0400", false},
		{"2019-09-12t13070400", true},
		{"2019-09-11t22-14-0400-2f4680648b319545c55c6149536f0a74527901f6", false},
	} {
		err := ValidateAppVersion(tc.in)
		if (err != nil) != tc.wantErr {
			t.Errorf("ValidateAppVersion(%q) = %v, want error = %t", tc.in, err, tc.wantErr)
		}
	}
}

func TestChooseOne(t *testing.T) {
	tests := []struct {
		configVar   string
		wantMatches string
	}{
		{"foo", "foo"},
		{"foo1 \n foo2", "^foo[12]$"},
		{"foo1\nfoo2", "^foo[12]$"},
		{"foo1 foo2", "^foo[12]$"},
	}
	for _, test := range tests {
		got, err := chooseOne(test.configVar)
		if err != nil {
			t.Fatal(err)
		}
		matched, err := regexp.MatchString(test.wantMatches, got)
		if err != nil {
			t.Fatal(err)
		}
		if !matched {
			t.Errorf("chooseOne(%q) = %q, _, want matches %q", test.configVar, got, test.wantMatches)
		}
	}
}

func TestProcessOverrides(t *testing.T) {
	tr := true
	f := false
	cfg := Config{
		DBHost: "origHost",
		DBName: "origName",
		Quota:  QuotaSettings{QPS: 1, Burst: 2, MaxEntries: 3, RecordOnly: &tr},
	}
	ov := `
        DBHost: newHost
        Quota:
           MaxEntries: 17
           RecordOnly: false
    `
	processOverrides(&cfg, []byte(ov))
	got := cfg
	want := Config{
		DBHost: "newHost",
		DBName: "origName",
		Quota:  QuotaSettings{QPS: 1, Burst: 2, MaxEntries: 17, RecordOnly: &f},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want, +got):\n%s", diff)
	}
}
