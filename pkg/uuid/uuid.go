package uuid

import "github.com/google/uuid"

func GenerateNameUUID(name string) string {
	return uuid.MustParse(name).String()
}

func GenerateUUID() string {
	return uuid.Must(uuid.NewRandom()).String()
}
