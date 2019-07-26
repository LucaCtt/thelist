package data

import "testing"

func p(x int) *int {
	return &x
}

func TestItem_IsValid(t *testing.T) {
	tests := []struct {
		name string
		i    *Item
		want bool
	}{
		{
			name: "ShowID is not nil",
			i: &Item{
				ShowID: p(1),
			},
			want: true,
		},
		{
			name: "ShowID is nil",
			i: &Item{
				ShowID: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.IsValid(); got != tt.want {
				t.Errorf("Item.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
