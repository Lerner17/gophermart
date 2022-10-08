package models

import "testing"

func TestRegisterUser_ValidatePassword(t *testing.T) {
	type fields struct {
		Login          string
		Password       string
		HashedPassword string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "empty password",
			fields: fields{
				Login:    "admin",
				Password: "",
			},
			want: false,
		},
		{
			name: "low password length",
			fields: fields{
				Login:    "admin",
				Password: "123",
			},
			want: false,
		},
		{
			name: "Only numeric password",
			fields: fields{
				Login:    "admin",
				Password: "123456789",
			},
			want: false,
		},
		{
			name: "password without uppercase letters",
			fields: fields{
				Login:    "admin",
				Password: "12@#x3456789",
			},
			want: false,
		},
		{
			name: "password without lower letters",
			fields: fields{
				Login:    "admin",
				Password: "12@#X3456789",
			},
			want: false,
		},
		{
			name: "password without spec letters",
			fields: fields{
				Login:    "admin",
				Password: "12xxaaX3456789",
			},
			want: false,
		},
		{
			name: "correct password",
			fields: fields{
				Login:    "admin",
				Password: "qaA!!3@1xka@21",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ru := &RegisterUser{
				Login:          tt.fields.Login,
				Password:       tt.fields.Password,
				HashedPassword: tt.fields.HashedPassword,
			}
			if got := ru.ValidatePassword(); got != tt.want {
				t.Errorf("RegisterUser.ValidatePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
