package buffer

import "testing"

func TestBuffer_InsertAt(t *testing.T) {
	tests := []struct {
		name     string
		initial  []BufRow
		at       Cursor
		content  [][]rune
		expected []BufRow
		wantErr  bool
	}{
		{
			name:    "insert in middle of line",
			initial: []BufRow{[]rune("hello")},
			at:      Cursor{X: 5, Y: 0},
			content: [][]rune{[]rune(" world")},
			expected: []BufRow{
				[]rune("hello world"),
			},
			wantErr: false,
		},
		{
			name:    "insert new line after first",
			initial: []BufRow{[]rune("foo")},
			at:      Cursor{X: 3, Y: 0},
			content: [][]rune{[]rune("bar"), []rune("baz")},
			expected: []BufRow{
				[]rune("foobar"),
				[]rune("baz"),
			},
			wantErr: false,
		},
		{
			name:    "invalid Y cursor",
			initial: []BufRow{[]rune("line")},
			at:      Cursor{X: 0, Y: 1},
			content: [][]rune{[]rune("oops")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Buffer{
				Rows:   append([]BufRow{}, tt.initial...),
				cursor: &Cursor{},
			}
			err := b.InsertAt(tt.at, tt.content)
			if (err != nil) != tt.wantErr {
				t.Fatalf("InsertAt() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if len(b.Rows) != len(tt.expected) {
					t.Fatalf("expected %d rows, got %d", len(tt.expected), len(b.Rows))
				}
				for i := range b.Rows {
					if string(b.Rows[i]) != string(tt.expected[i]) {
						t.Errorf("row %v mismatch: got %v, want %v", i, b.Rows[i], tt.expected[i])
					}
				}
			}
		})
	}
}


func TestBuffer_DeleteAt(t *testing.T) {
	tests := []struct {
		name        string
		initial     []BufRow
		start, end  Cursor
		want        [][]rune
		wantRemain  []BufRow
		expectError bool
	}{
		{
			name:    "single-line delete",
			initial: []BufRow{[]rune("hello world")},
			start:   Cursor{X: 5, Y: 0},
			end:     Cursor{X: 11, Y: 0},
			want:    [][]rune{[]rune(" world")},
			wantRemain: []BufRow{
				[]rune("hello"),
			},
			expectError: false,
		},
		{
			name: "multi-line delete",
			initial: []BufRow{
				[]rune("line1"),
				[]rune("line2"),
				[]rune("line3"),
			},
			start:  Cursor{X: 2, Y: 0},
			end:    Cursor{X: 2, Y: 2},
			want:   [][]rune{[]rune("ne1"), []rune("line2"), []rune("lin")},
			wantRemain: []BufRow{
				[]rune("lie3"),
			},
			expectError: false,
		},
		{
			name:        "invalid end cursor",
			initial:     []BufRow{[]rune("only line")},
			start:       Cursor{X: 0, Y: 0},
			end:         Cursor{X: 0, Y: 1},
			expectError: true,
		},
		{
			name:        "start > end",
			initial:     []BufRow{[]rune("bad")},
			start:       Cursor{X: 0, Y: 1},
			end:         Cursor{X: 0, Y: 0},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Buffer{
				Rows: append([]BufRow{}, tt.initial...),
			}
			got, err := b.DeleteAt(tt.start, tt.end)
			if (err != nil) != tt.expectError {
				t.Fatalf("DeleteAt() error = %v, wantErr %v", err, tt.expectError)
			}
			if !tt.expectError {
				// check deleted content
				for i := range got {
					if string(got[i]) != string(tt.want[i]) {
						t.Errorf("deleted line %d = %q, want %q", i, got[i], tt.want[i])
					}
				}
				// check remaining buffer
				for i := range b.Rows {
					if string(b.Rows[i]) != string(tt.wantRemain[i]) {
						t.Errorf("remaining line %d = %q, want %q", i, b.Rows[i], tt.wantRemain[i])
					}
				}
			}
		})
	}
}

