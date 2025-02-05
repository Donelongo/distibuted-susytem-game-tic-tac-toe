package utils_test

import (
	"go-xox-grpc-ai/internal/utils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to simulate user input
func mockInput(input string) func() {
	// Create a pipe to simulate stdin
	r, w, _ := os.Pipe()
	os.Stdin = r

	// Write input and close the writer to simulate user pressing Enter
	w.WriteString(input + "\n")
	w.Close()

	// Return a function to restore stdin
	return func() {
		os.Stdin = os.NewFile(uintptr(0), "/dev/stdin")
	}
}

func TestGetUserInput_ValidInput(t *testing.T) {
	restore := mockInput("1")
	defer restore() // Restore original stdin after test

	// Call function with valid choices
	result := utils.GetUserInput("1", "2", "3", "X", "O")

	// Assert result matches expected value
	assert.Equal(t, "1", result, "Expected input to be '1'")
}

func TestGetUserInput_InvalidThenValidInput(t *testing.T) {
	restore := mockInput("invalid\n2")
	defer restore()

	// Call function, ensuring it only accepts valid choices
	result := utils.GetUserInput("1", "2", "3", "X", "O")

	// Assert result is the valid second input
	assert.Equal(t, "2", result, "Expected input to be '2'")
}

func TestGetUserInput_CaseInsensitive(t *testing.T) {
	restore := mockInput("x")
	defer restore()

	// Call function with case-insensitive check
	result := utils.GetUserInput("1", "2", "3", "X", "O")

	// Assert result is uppercase as expected
	assert.Equal(t, "X", result, "Expected 'x' to be converted to 'X'")
}
