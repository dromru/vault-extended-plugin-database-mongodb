package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// errCodeRoleAlreadyExists is the numeric code MongoDB returns from
// insertRoleDocument when a createRole command collides with an existing role
// (see user_management_commands.cpp in the MongoDB source). The code is not
// registered in the central error_codes.yml, so the wire codeName is the
// generic "Location51002" and must not be matched by name.
const errCodeRoleAlreadyExists = 51002

// ensureRoles creates each custom role if it does not already exist. Custom
// roles are shared across users, so a role already created for a previous user
// is left untouched.
func (m *MongoDB) ensureRoles(ctx context.Context, db string, roles []mongoCustomRole) error {
	if len(roles) == 0 {
		return nil
	}

	client, err := m.Connection(ctx)
	if err != nil {
		return err
	}

	for _, role := range roles {
		if err := m.ensureRole(ctx, client, db, role); err != nil {
			return err
		}
	}

	return nil
}

// ensureRole creates the role if it does not already exist.
func (m *MongoDB) ensureRole(ctx context.Context, client *mongo.Client, db string, role mongoCustomRole) error {
	exists, err := roleExists(ctx, client, db, role.Role)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	cmd := &createRoleCommand{
		Role:       role.Role,
		Privileges: role.Privileges,
		Roles:      []string{},
	}
	if err := m.runCommandWithRetry(ctx, db, cmd); err != nil {
		if isRoleAlreadyExists(err) {
			return nil
		}
		return err
	}

	return nil
}

func roleExists(ctx context.Context, client *mongo.Client, db, role string) (bool, error) {
	var out struct {
		Roles []bson.M `bson:"roles"`
	}
	err := client.Database(db).RunCommand(ctx, bson.D{
		{Key: "rolesInfo", Value: []string{role}},
		{Key: "showPrivileges", Value: false},
	}).Decode(&out)
	if err != nil {
		return false, err
	}
	return len(out.Roles) > 0, nil
}

func isRoleAlreadyExists(err error) bool {
	var cmdErr mongo.CommandError
	if !errors.As(err, &cmdErr) {
		return false
	}
	return cmdErr.Code == errCodeRoleAlreadyExists
}
