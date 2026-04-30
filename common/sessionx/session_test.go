package sessionx

import "testing"

func TestManagerEncodeDecode(t *testing.T) {
	manager := NewManager(Config{
		Name:   "test-session",
		Secret: "test-secret",
		MaxAge: 60,
	})

	encoded, err := manager.Encode(42)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	if encoded == "" {
		t.Fatal("Encode() returned empty session")
	}

	userID, err := manager.Decode(encoded)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if userID != 42 {
		t.Fatalf("Decode() userID = %d, want 42", userID)
	}
}
