package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabets = "abcdefghijklmnopqrstuvwxyx"

// RandomInt generates a random int b/w 0 and max - min
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates random string of n characters
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabets)

	for i := 0; i < n; i++ {
		c := alphabets[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// Generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// generate random money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

// generates a random product title
func RandomProductTitle() string {
	return fmt.Sprintf("Product %s", RandomString(5))
}

// generates a random product description
func RandomDescription() string {
	return fmt.Sprintf("This is a description about %s.", RandomString(10))
}

// Returns a random price between 5.00 and 25.00 as float64
func RandomPrice() float64 {
	return float64(rand.Intn(2000)+500) / 100.0 // e.g., 5.00 to 25.00
}

// returns a random categoryID
func RandomCategoryID() int64 {
	return RandomInt(1, 5)
}

// creates a random status (completed or pending)
func RandomStatus() string {
	statuses := []string{"completed", "pending"}
	return statuses[rand.Intn(len(statuses))]
}

// creates a random pay method (stripe or payoneer)
func RandomPaymentMethod() string {
	methods := []string{"stripe", "payoneer"}
	randomValue := methods[rand.Intn(len(methods))]

	return randomValue
}

// gives random rating b/w 1 and 5
func RandomRating() int64 {
	return int64(rand.Intn(5) + 1) // gives 1 to 5
}

// random comment from 20 to 39 characters
func RandomComment() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ ")

	length := rand.Intn(20) + 20 // 20 to 39 characters
	comment := make([]rune, length)
	for i := range comment {
		comment[i] = letters[rand.Intn(len(letters))]
	}
	return string(comment)
}
