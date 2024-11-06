package app

import (
	"math/rand"
)

func genData(numRecords int) []User {
	users := make([]User, numRecords)
	for i := 0; i < numRecords; i++ {
		users[i] = User{
			Username: generateRandomString(7),
			Password: generateRandomString(7),
			City:     generateRandomString(5),
		}
	}
	return users
}

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
