package captain_auth_handler

import (
	"testing"
)

func TestLogin(t *testing.T) {
	// os.Setenv("POSTGRES_DB", "rapido")
	// ctx := context.Background()
	// svc, err := initAuthHandler()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// // 1. Login (Generate OTP) - Verifies DB and Redis connection
	// data, err := svc.Login(ctx, &dto.LoginRequest{
	// 	Phone: "+919876543210",
	// })
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println("Login successful:", data.Message)
}

// TestRefresh requires mocking repo to intercept OTP or verify without valid OTP.
// Skipping full flow test here. Use manual verification with /captain/auth/refresh.
