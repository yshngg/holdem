package player

import (
	"encoding/base64"
	"testing"

	"github.com/google/uuid"
)

func TestUUID(t *testing.T) {
	id := uuid.New()
	name := base64.StdEncoding.EncodeToString([]byte(id.String()))[:7]
	t.Logf("id: %s, name: %s", id, name)
}
