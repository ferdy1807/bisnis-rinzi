package utils

import (
	"crypto/rand"
	"fmt"
)

// GenerateUUIDv4 melahirkan string UUID v4 murni standar RFC 4122
// yang valid dan aman untuk menduduki Primary Key PostgreSQL.
func GenerateUUIDv4() string {
	uuid := make([]byte, 16)
	_, _ = rand.Read(uuid)

	// Set varian ke RFC 4122
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	// Set versi ke 4 (random)
	uuid[6] = (uuid[6] & 0x0f) | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:],
	)
}
