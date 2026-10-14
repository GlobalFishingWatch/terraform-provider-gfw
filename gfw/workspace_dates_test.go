package gfw

import (
	"encoding/json"
	"testing"

	"github.com/globalfishingwatch.org/terraform-provider-gfw/gfw/api"
)

func TestWorkspaceDatesPatchBody(t *testing.T) {
	cases := map[string]struct {
		ws   api.CreateWorkspace
		want string
	}{
		"unchanged is omitted": {api.CreateWorkspace{}, `{}`},
		"removed is null":      {api.CreateWorkspace{StartAt: jsonStringOrNull(""), EndAt: jsonStringOrNull("")}, `{"startAt":null,"endAt":null}`},
		"set is string":        {api.CreateWorkspace{StartAt: jsonStringOrNull("2024-01-01T00:00:00.000Z")}, `{"startAt":"2024-01-01T00:00:00.000Z"}`},
	}
	for name, c := range cases {
		b, err := json.Marshal(c.ws)
		if err != nil || string(b) != c.want {
			t.Errorf("%s: got %s (err %v), want %s", name, b, err, c.want)
		}
	}
}
