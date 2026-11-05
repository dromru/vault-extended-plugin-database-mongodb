package mongodb

import (
	"encoding/json"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestMongoDBStatementUnmarshal locks the creation-statement contract between
// Protector and the plugin.
func TestMongoDBStatementUnmarshal(t *testing.T) {
	const statement = `{
		"db": "admin",
		"roles": [
			{"role": "readWrite", "db": "orders"},
			{"role": "readAnyDatabase"},
			{"role": "protector_ab12cd34ef56", "db": "admin"}
		],
		"custom_roles": [
			{
				"role": "protector_ab12cd34ef56",
				"privileges": [
					{"resource": {"db": "orders", "collection": "invoices"}, "actions": ["find", "insert", "update", "remove"]}
				]
			}
		],
		"authentication_restrictions": [{"clientSource": ["10.0.0.0/16", "10.1.2.3"]}]
	}`

	var cs mongoDBStatement
	if err := json.Unmarshal([]byte(statement), &cs); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if cs.DB != "admin" {
		t.Errorf("DB = %q, want admin", cs.DB)
	}
	if len(cs.Roles) != 3 {
		t.Fatalf("len(Roles) = %d, want 3", len(cs.Roles))
	}
	if len(cs.CustomRoles) != 1 {
		t.Fatalf("len(CustomRoles) = %d, want 1", len(cs.CustomRoles))
	}

	er := cs.CustomRoles[0]
	if er.Role != "protector_ab12cd34ef56" {
		t.Errorf("CustomRoles[0] = %+v", er)
	}
	if len(er.Privileges) != 1 {
		t.Fatalf("len(Privileges) = %d, want 1", len(er.Privileges))
	}
	p := er.Privileges[0]
	if p.Resource.DB != "orders" || p.Resource.Collection != "invoices" {
		t.Errorf("Resource = %+v", p.Resource)
	}
	if len(p.Actions) != 4 {
		t.Errorf("Actions = %v", p.Actions)
	}
	if len(cs.AuthenticationRestrictions) != 1 ||
		len(cs.AuthenticationRestrictions[0].ClientSource) != 2 {
		t.Errorf("AuthenticationRestrictions = %+v", cs.AuthenticationRestrictions)
	}
}

// TestCreateRoleCommandBSON verifies the createRole command marshals to the
// shape MongoDB expects.
func TestCreateRoleCommandBSON(t *testing.T) {
	cmd := &createRoleCommand{
		Role: "protector_ab12cd34ef56",
		Privileges: []mongoPrivilege{
			{Resource: mongoResource{DB: "orders", Collection: "invoices"}, Actions: []string{"find"}},
		},
		Roles: []string{},
	}
	doc, err := bson.Marshal(cmd)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m bson.M
	if err := bson.Unmarshal(doc, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := m["createRole"]; !ok {
		t.Errorf("expected createRole key, got %v", m)
	}
	if _, ok := m["privileges"]; !ok {
		t.Errorf("expected privileges key, got %v", m)
	}
	roles, ok := m["roles"].(primitive.A)
	if !ok {
		t.Fatalf("expected roles key with array value, got %T: %v", m["roles"], m["roles"])
	}
	if len(roles) != 0 {
		t.Errorf("roles = %v, want empty array", roles)
	}
}

// TestMongoResourceBSON verifies that a database-level privilege keeps the
// empty collection field, as required by the createRole/createUser commands.
func TestMongoResourceBSON(t *testing.T) {
	priv := mongoPrivilege{
		Resource: mongoResource{DB: "admin"},
		Actions:  []string{"find"},
	}
	doc, err := bson.Marshal(priv)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var m bson.M
	if err := bson.Unmarshal(doc, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	res, ok := m["resource"].(bson.M)
	if !ok {
		t.Fatalf("expected resource object, got %T: %v", m["resource"], m["resource"])
	}
	if _, present := res["collection"]; !present {
		t.Errorf("resource = %v, want collection key present (empty string)", res)
	}
	if res["collection"] != "" {
		t.Errorf("collection = %v, want empty string", res["collection"])
	}
}

// TestToStandardRolesArray verifies mixed string/document roles are preserved.
func TestToStandardRolesArray(t *testing.T) {
	roles := mongodbRoles{
		{Role: "readAnyDatabase"},
		{Role: "readWrite", DB: "orders"},
	}
	got := roles.toStandardRolesArray()
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if s, ok := got[0].(string); !ok || s != "readAnyDatabase" {
		t.Errorf("got[0] = %v (%T), want string readAnyDatabase", got[0], got[0])
	}
	if _, ok := got[1].(mongodbRole); !ok {
		t.Errorf("got[1] type = %T, want mongodbRole", got[1])
	}
}
