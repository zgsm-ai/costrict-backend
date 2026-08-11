package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/zgsm-ai/chat-rag/internal/types"
)

// mockTransport implements RoundTripper interface
type mockTransport struct{}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Read request body to determine test case
	body, _ := io.ReadAll(req.Body)
	var request struct {
		Messages []types.Message `json:"messages"`
	}
	json.Unmarshal(body, &request)

	resp := &http.Response{
		Header: make(http.Header),
	}

	// Handle different test cases
	switch len(request.Messages) {
	case 0: // Empty messages case
		resp.StatusCode = http.StatusBadRequest
		resp.Body = io.NopCloser(bytes.NewBufferString(`{
			"error": "messages cannot be empty"
		}`))
	default: // All other cases
		resp.StatusCode = http.StatusOK
		resp.Body = io.NopCloser(bytes.NewBufferString(`{
			"choices": [{
				"message": {
					"content": "mocked-summary"
				}
			}]
		}`))
	}

	return resp, nil
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestLLMClient_NonStreamingResponseMetadata(t *testing.T) {
	headers := make(http.Header)
	client := &LLMClient{
		modelName: "test-model",
		endpoint:  "http://mock-endpoint/v1/chat/completions",
		headers:   &headers,
		httpClient: &http.Client{Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
			responseHeaders := make(http.Header)
			responseHeaders.Set("X-Custom-Trace-ID", "trace-success")
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     responseHeaders,
				Body: io.NopCloser(bytes.NewBufferString(`{
					"id":"response-id",
					"choices":[],
					"usage":{}
				}`)),
			}, nil
		})},
	}

	result, err := client.ChatLLMWithMessagesRaw(context.Background(), types.LLMRequestParams{}, nil)
	if err != nil {
		t.Fatalf("ChatLLMWithMessagesRaw() error = %v", err)
	}
	if result.StatusCode != http.StatusOK {
		t.Fatalf("StatusCode = %d, want %d", result.StatusCode, http.StatusOK)
	}
	if got := result.Header.Get("x-custom-trace-id"); got != "trace-success" {
		t.Fatalf("trace header = %q, want %q", got, "trace-success")
	}
	if result.Response.Id != "response-id" {
		t.Fatalf("response ID = %q, want %q", result.Response.Id, "response-id")
	}
}

func TestLLMClient_NonStreamingErrorResponseMetadata(t *testing.T) {
	headers := make(http.Header)
	client := &LLMClient{
		modelName: "test-model",
		endpoint:  "http://mock-endpoint/v1/chat/completions",
		headers:   &headers,
		httpClient: &http.Client{Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
			responseHeaders := make(http.Header)
			responseHeaders.Set("X-Custom-Trace-ID", "trace-error")
			return &http.Response{
				StatusCode: http.StatusBadGateway,
				Header:     responseHeaders,
				Body:       io.NopCloser(bytes.NewBufferString(`{"error":{"message":"upstream failed"}}`)),
			}, nil
		})},
	}

	result, err := client.ChatLLMWithMessagesRaw(context.Background(), types.LLMRequestParams{}, nil)
	if err == nil {
		t.Fatal("ChatLLMWithMessagesRaw() error = nil, want upstream error")
	}
	if result.StatusCode != http.StatusBadGateway {
		t.Fatalf("StatusCode = %d, want %d", result.StatusCode, http.StatusBadGateway)
	}
	if got := result.Header.Get("x-custom-trace-id"); got != "trace-error" {
		t.Fatalf("trace header = %q, want %q", got, "trace-error")
	}
}

func TestLLMClient_StreamingErrorResponseHeaders(t *testing.T) {
	headers := make(http.Header)
	client := &LLMClient{
		modelName: "test-model",
		endpoint:  "http://mock-endpoint/v1/chat/completions",
		headers:   &headers,
		httpClient: &http.Client{Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
			responseHeaders := make(http.Header)
			responseHeaders.Set("X-Custom-Trace-ID", "trace-error")
			return &http.Response{
				StatusCode: http.StatusBadGateway,
				Header:     responseHeaders,
				Body:       io.NopCloser(bytes.NewBufferString(`{"error":{"message":"upstream failed"}}`)),
			}, nil
		})},
	}

	var callbackResponse LLMResponse
	err := client.ChatLLMWithMessagesStreamRaw(
		context.Background(),
		types.LLMRequestParams{},
		nil,
		func(response LLMResponse) error {
			callbackResponse = response
			return nil
		},
	)
	if err == nil {
		t.Fatal("ChatLLMWithMessagesStreamRaw() error = nil, want upstream error")
	}
	if !callbackResponse.HeaderOnly {
		t.Fatal("expected a header-only callback before the upstream error")
	}
	if callbackResponse.StatusCode != http.StatusBadGateway {
		t.Fatalf("StatusCode = %d, want %d", callbackResponse.StatusCode, http.StatusBadGateway)
	}
	if got := callbackResponse.Header.Get("x-custom-trace-id"); got != "trace-error" {
		t.Fatalf("trace header = %q, want %q", got, "trace-error")
	}
}

func TestLLMClient_ChatLLMWithMessages_FormatCheck(t *testing.T) {
	// Simple message format validation without creating an actual client
	messages := []struct {
		Role    string
		Content string
	}{
		{"system", "You are a helpful assistant that summarizes content."},
		{"user", "This is test content that needs to be summarized"},
	}

	// Verify messages contain expected content
	for _, msg := range messages {
		if msg.Content == "" {
			t.Errorf("Message content should not be empty")
		}
		if msg.Role == "" {
			t.Errorf("Message role should not be empty")
		}
	}
}

func TestLLMClient_ChatLLMWithMessages(t *testing.T) {
	// Setup mock HTTP transport
	mockTransport := &mockTransport{}

	// Create test client with mock transport
	headers := make(http.Header)
	headers.Add("Content-Type", "application/json")

	client := &LLMClient{
		modelName:  "test-model",
		endpoint:   "http://mock-endpoint/v1/chat/completions",
		httpClient: &http.Client{Transport: mockTransport},
		headers:    &headers,
	}

	// Test cases
	testCases := []struct {
		name        string
		messages    []types.Message
		expectError bool
	}{
		{
			name:        "Empty messages",
			messages:    []types.Message{},
			expectError: false, // Empty messages are not an error
		},
		{
			name: "Single user message",
			messages: []types.Message{
				{Role: "user", Content: "This is a short text to summarize."},
			},
			expectError: false,
		},
		{
			name: "System and user messages",
			messages: []types.Message{
				{Role: "system", Content: "You are a helpful assistant that summarizes content."},
				{Role: "user", Content: "This is a longer text that contains multiple sentences. It discusses various topics and should be summarized properly. The summary should retain the key information while being concise."},
			},
			expectError: false,
		},
		{
			name: "Conversation with assistant",
			messages: []types.Message{
				{Role: "system", Content: "You are a helpful assistant."},
				{Role: "user", Content: "Please summarize this content."},
				{Role: "assistant", Content: "I'll help you summarize the content."},
				{Role: "user", Content: "Here is the content: This is important information that needs to be condensed."},
			},
			expectError: false,
		},
	}

	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			summary, err := client.GenerateContent(ctx, "", tc.messages)

			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if err == nil {
				// Check that summary is not empty when we successfully call the API
				if summary == "" {
					t.Error("Received empty summary")
				}
			}
			fmt.Println("messages:", tc.messages)
			fmt.Println("summary:", summary)
		})
	}
}

func TestLLMClient_Integration(t *testing.T) {
	// Skip this test as it requires external services
	t.Skip("Skipping integration test")

	// Note: This test is skipped. The following code is just an example.
	// If you want to run this test, remove t.Skip() and add the following import:
	// import "context"
}
