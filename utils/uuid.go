package utils

import (
	"strings"

	"github.com/google/uuid"
)

// e83ca1359e7c4266b0f18def28611681
func UUID() string {
	id := uuid.New()
	idWithoutHyphen := id.String()
	return strings.ReplaceAll(idWithoutHyphen, "-", "")
}
