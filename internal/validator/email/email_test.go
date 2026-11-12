package email

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr error
	}{
		{"email válido", "user@example.com", nil},
		{"email vazio", "", ErrEmailEmpty},
		{"email vazio com espaços", "   ", ErrEmailEmpty},
		{"sem @", "userexample.com", ErrEmailInvalid},
		{"sem domínio", "user@", ErrEmailInvalid},
		{"sem TLD", "user@localhost", ErrEmailInvalid},
		{"com espaço no meio", "user @example.com", ErrEmailInvalid},
		{"com nome junto", "User <user@example.com>", ErrEmailInvalid},
		{"subdomínio válido", "user@mail.example.com", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if err != tt.wantErr {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tt.email, err, tt.wantErr)
			}
		})
	}
}
