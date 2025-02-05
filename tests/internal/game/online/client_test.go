package online_test

import (
	"context"
	"fmt"
	"go-xox-grpc-ai/internal/game"
	"go-xox-grpc-ai/internal/game/online"
	"go-xox-grpc-ai/internal/game/online/api"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockGameServiceClient is a mock of the GameServiceClient interface
type MockGameServiceClient struct {
	mock.Mock
}

func (m *MockGameServiceClient) Join(ctx context.Context, in *api.JoinRequest, opts ...grpc.CallOption) (*api.JoinResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.JoinResponse), args.Error(1)
}

func (m *MockGameServiceClient) ServerMove(ctx context.Context, in *api.ServerMoveRequest, opts ...grpc.CallOption) (api.GameService_ServerMoveClient, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(api.GameService_ServerMoveClient), args.Error(1)
}

func (m *MockGameServiceClient) ClientMove(ctx context.Context, in *api.ClientMoveRequest, opts ...grpc.CallOption) (*api.ClientMoveResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.ClientMoveResponse), args.Error(1)
}

// Helper function to set unexported fields using reflection
func setUnexportedField(i interface{}, fieldName string, value interface{}) error {
	v := reflect.ValueOf(i).Elem() // get the value of the struct
	f := v.FieldByName(fieldName)
	if !f.IsValid() {
		return fmt.Errorf("no such field: %s", fieldName)
	}
	if !f.CanSet() {
		return fmt.Errorf("cannot set field: %s", fieldName)
	}
	f.Set(reflect.ValueOf(value))
	return nil
}

// TestJoinGame verifies that the client can join the game
func TestJoinGame(t *testing.T) {
	// Create a mock gRPC client
	mockGameServiceClient := new(MockGameServiceClient)

	// Initialize the Client
	client := &online.Client{}

	// Set the mock gRPC client using reflection
	err := setUnexportedField(client, "grpcClient", mockGameServiceClient)
	if err != nil {
		t.Fatalf("failed to set grpcClient: %v", err)
	}

	// Mock the Join response
	mockJoinResponse := &api.JoinResponse{
		Success:      true,
		ClientPlayer: game.PLAYER_X,
	}
	mockGameServiceClient.On("Join", mock.Anything, mock.Anything).Return(mockJoinResponse, nil)

	// Call the JoinGame method
	client.JoinGame("localhost:5000")

	// Verify that the Join method was called and expectations were met
	mockGameServiceClient.AssertExpectations(t)
}

// TestClientMoves verifies the behavior of the ClientMoves method
func TestClientMoves(t *testing.T) {
	// Create a mock gRPC client
	mockGameServiceClient := new(MockGameServiceClient)

	// Initialize the Client
	client := &online.Client{}

	// Set the mock gRPC client using reflection
	err := setUnexportedField(client, "grpcClient", mockGameServiceClient)
	if err != nil {
		t.Fatalf("failed to set grpcClient: %v", err)
	}

	// Mock the Join response
	mockJoinResponse := &api.JoinResponse{
		Success:      true,
		ClientPlayer: game.PLAYER_X,
	}
	mockGameServiceClient.On("Join", mock.Anything, mock.Anything).Return(mockJoinResponse, nil)

	// Mock the ClientMove response
	mockClientMoveResponse := &api.ClientMoveResponse{
		Success:        true,
		Board:          []string{"X", "", "", "", "", "", "", "", ""},
		IsGameFinished: false,
	}
	mockGameServiceClient.On("ClientMove", mock.Anything, mock.Anything).Return(mockClientMoveResponse, nil)

	// Call the JoinGame method
	client.JoinGame("localhost:5000")

	// Simulate a move
	err = client.ClientMoves(context.Background())

	// Verify the ClientMove method was called and expectations were met
	mockGameServiceClient.AssertExpectations(t)

	// Assert no error occurred during the ClientMoves process
	assert.NoError(t, err)
}

// TestMovePlayed verifies that the MovePlayed method updates the board correctly
func TestMovePlayed(t *testing.T) {
	// Create a mock gRPC client
	mockGameServiceClient := new(MockGameServiceClient)

	// Initialize the Client
	client := &online.Client{}

	// Set the mock gRPC client using reflection
	err := setUnexportedField(client, "grpcClient", mockGameServiceClient)
	if err != nil {
		t.Fatalf("failed to set grpcClient: %v", err)
	}

	// Mock the Join response
	mockJoinResponse := &api.JoinResponse{
		Success:      true,
		ClientPlayer: game.PLAYER_X,
	}
	mockGameServiceClient.On("Join", mock.Anything, mock.Anything).Return(mockJoinResponse, nil)

	// Mock the ClientMove response
	mockClientMoveResponse := &api.ClientMoveResponse{
		Success:        true,
		Board:          []string{"X", "O", "X", "", "", "", "", "", ""},
		IsGameFinished: false,
	}
	mockGameServiceClient.On("ClientMove", mock.Anything, mock.Anything).Return(mockClientMoveResponse, nil)

	// Call the JoinGame method
	client.JoinGame("localhost:5000")

	// Simulate the move being played
	client.MovePlayed(mockClientMoveResponse.Board, mockClientMoveResponse.IsGameFinished)

	// Assert the board state has been updated
	mockGameServiceClient.AssertExpectations(t)
}
