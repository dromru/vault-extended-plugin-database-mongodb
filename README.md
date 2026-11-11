# Vault database secrets engine - Extended MongoDB plugin

This plugin is a modification of the [official Vault plugin](https://github.com/hashicorp/vault/tree/main/plugins/database/mongodb).

## What was changed

1. Added the ability to pass privileges and create roles based on them. This makes it possible to specify privileges on collections.
2. Added the ability to specify the user's allowed IPs.

## Limitations

- Custom roles `<serviceName>_<hash>` are created in the admin database once per unique
  `(database, collection, actions)` combination and are shared by all users.
  When a user is deleted (revoke), the roles are NOT removed — they remain in MongoDB
  ("orphaned" roles) and accumulate as new privilege combinations appear.
  Cleanup is done manually: roles not assigned to any user
  can be found via `rolesInfo` and dropped with the `dropRole` command.
- The connection user (connection_url) must have the right to create roles
  and read `rolesInfo` in the admin database (for example, the `root` or `clusterAdmin` role).
