package main

import "testing"

func TestSNSRegionFromARN(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		arn     string
		want    string
		wantErr bool
	}{
		{
			name: "valid arn extracts region",
			arn:  "arn:aws:sns:eu-north-1:123456789012:book-created",
			want: "eu-north-1",
		},
		{
			name:    "invalid arn fails",
			arn:     "not-an-arn",
			wantErr: true,
		},
		{
			name:    "missing region fails",
			arn:     "arn:aws:sns::123456789012:book-created",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := snsRegionFromARN(tc.arn)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for arn %q", tc.arn)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}
