// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package mongodb

import "go.mongodb.org/mongo-driver/mongo/writeconcern"

type createUserCommand struct {
	Username                   string                             `bson:"createUser"`
	Password                   string                             `bson:"pwd,omitempty"`
	Roles                      []interface{}                      `bson:"roles"`
	AuthenticationRestrictions []*mongoAuthenticationRestrictions `bson:"authenticationRestrictions,omitempty"`
}

type updateUserCommand struct {
	Username string `bson:"updateUser"`
	Password string `bson:"pwd"`
}

type dropUserCommand struct {
	Username     string                     `bson:"dropUser"`
	WriteConcern *writeconcern.WriteConcern `bson:"writeConcern"`
}

type mongoResource struct {
	DB         string `json:"db" bson:"db"`
	Collection string `json:"collection,omitempty" bson:"collection"`
}

type mongoPrivilege struct {
	Resource mongoResource `json:"resource" bson:"resource"`
	Actions  []string      `json:"actions" bson:"actions"`
}

type mongoCustomRole struct {
	Role       string           `json:"role" bson:"role"`
	Privileges []mongoPrivilege `json:"privileges" bson:"privileges"`
}

type mongoAuthenticationRestrictions struct {
	ClientSource  []string `json:"clientSource,omitempty" bson:"clientSource,omitempty"`
	ServerAddress []string `json:"serverAddress,omitempty" bson:"serverAddress,omitempty"`
}

type createRoleCommand struct {
	Role       string           `bson:"createRole"`
	Privileges []mongoPrivilege `bson:"privileges"`
	Roles      []string         `bson:"roles"`
}

type mongodbRole struct {
	Role string `json:"role" bson:"role"`
	DB   string `json:"db"   bson:"db"`
}

type mongodbRoles []mongodbRole

type mongoDBStatement struct {
	DB                         string                             `json:"db"`
	Roles                      mongodbRoles                       `json:"roles"`
	CustomRoles                []mongoCustomRole                  `json:"custom_roles"`
	AuthenticationRestrictions []*mongoAuthenticationRestrictions `json:"authentication_restrictions"`
}

// Convert array of role documents like:
//
// [ { "role": "readWrite" }, { "role": "readWrite", "db": "test" } ]
//
// into a "standard" MongoDB roles array containing both strings and role documents:
//
// [ "readWrite", { "role": "readWrite", "db": "test" } ]
//
// MongoDB's createUser command accepts the latter.
func (roles mongodbRoles) toStandardRolesArray() []interface{} {
	var standardRolesArray []interface{}
	for _, role := range roles {
		if role.DB == "" {
			standardRolesArray = append(standardRolesArray, role.Role)
		} else {
			standardRolesArray = append(standardRolesArray, role)
		}
	}
	return standardRolesArray
}
