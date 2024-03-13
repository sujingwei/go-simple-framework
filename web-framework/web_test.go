package webframework

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestUUID(t *testing.T) {
	uid := strings.ReplaceAll(uuid.New().String(), "-", "")
	fmt.Printf("%T, %v, %s\n", uid, uid, uid)
}
