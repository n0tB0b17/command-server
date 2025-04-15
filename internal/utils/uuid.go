package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func GenerateUUID(id_type string) string {
	id := uuid.New().String()[0:8]
	t := time.Now().Format("20060102150405")
	newID := fmt.Sprintf("%s_%s_%s", id_type, t, id)

	return newID
}
