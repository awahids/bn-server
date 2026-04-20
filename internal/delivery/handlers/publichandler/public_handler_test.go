package publichandler

import "testing"

func TestIsAllowedAudioHost_ExactMatchOnly(t *testing.T) {
	original := allowedAudioDomains
	allowedAudioDomains = []string{"localhost", "127.0.0.1"}
	defer func() { allowedAudioDomains = original }()

	tests := []struct {
		name string
		host string
		want bool
	}{
		{name: "localhost exact", host: "localhost", want: true},
		{name: "localhost uppercase", host: "LOCALHOST", want: true},
		{name: "loopback exact", host: "127.0.0.1", want: true},
		{name: "localhost suffix attack", host: "localhost.attacker.com", want: false},
		{name: "loopback suffix attack", host: "127.0.0.1.attacker.com", want: false},
		{name: "empty host", host: "", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isAllowedAudioHost(tt.host)
			if got != tt.want {
				t.Fatalf("isAllowedAudioHost(%q) = %v, want %v", tt.host, got, tt.want)
			}
		})
	}
}

func TestIsAllowedAudioHost_WildcardRule(t *testing.T) {
	original := allowedAudioDomains
	allowedAudioDomains = []string{"*.example.com"}
	defer func() { allowedAudioDomains = original }()

	tests := []struct {
		host string
		want bool
	}{
		{host: "example.com", want: true},
		{host: "api.example.com", want: true},
		{host: "deep.api.example.com", want: true},
		{host: "badexample.com", want: false},
	}

	for _, tt := range tests {
		got := isAllowedAudioHost(tt.host)
		if got != tt.want {
			t.Fatalf("isAllowedAudioHost(%q) = %v, want %v", tt.host, got, tt.want)
		}
	}
}
