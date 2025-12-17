package guard

import (
	"encoding/json"
	"reflect"
	"testing"
)

// 测试函数
func TestGetModuleNameInline(t *testing.T) {
	raw := []byte(`{
		"name": "egress.endpoint",
		"key": "endpoint",
		"cfg": "xxxxxxx"
	}`)

	name, cfg, err := getModuleNameInline("name", raw)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	t.Logf("module name: %s", name)
	t.Logf("remaining config: %s", string(cfg))

	if name != "egress.endpoint" {
		t.Errorf("expected module name 'egress.endpoint', got %q", name)
	}

	var m map[string]any
	_ = json.Unmarshal(cfg, &m)
	if _, ok := m["name"]; ok {
		t.Errorf("expected 'name' key to be removed, but still found it")
	}
}

func TestParseStructTag(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      map[string]string
		expectErr bool
	}{
		{
			name:  "normal case",
			input: `key1=val1 key2=val2`,
			want: map[string]string{
				"key1": "val1",
				"key2": "val2",
			},
			expectErr: false,
		},
		{
			name:  "extra spaces",
			input: ` key1=val1   key2=val2 `,
			want: map[string]string{
				"key1": "val1",
				"key2": "val2",
			},
			expectErr: false,
		},
		{
			name:      "missing equal",
			input:     `key1 val2=val2`,
			want:      nil,
			expectErr: true,
		},
		{
			name:  "empty string",
			input: ``,
			want:  map[string]string{},
		},
		{
			name:  "single key",
			input: `foo=bar`,
			want: map[string]string{
				"foo": "bar",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStructTag(tt.input)
			if (err != nil) != tt.expectErr {
				t.Fatalf("ParseStructTag() error = %v, expectErr %v", err, tt.expectErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseStructTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
