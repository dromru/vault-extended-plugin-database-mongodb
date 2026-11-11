package mongodb

import (
	"errors"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestIsRoleAlreadyExists(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "code 51002 with Location codeName",
			err:  mongo.CommandError{Code: 51002, Name: "Location51002", Message: `Role "x@admin" already exists`},
			want: true,
		},
		{
			name: "wrapped code 51002",
			err:  fmt.Errorf("run command: %w", mongo.CommandError{Code: 51002, Name: "Location51002"}),
			want: true,
		},
		{
			name: "other code",
			err:  mongo.CommandError{Code: 40414, Name: "IDLFailedToParse", Message: "BSON field 'createRole.roles' is missing"},
			want: false,
		},
		{
			name: "unrelated error",
			err:  errors.New("connection refused"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRoleAlreadyExists(tt.err); got != tt.want {
				t.Errorf("isRoleAlreadyExists(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
