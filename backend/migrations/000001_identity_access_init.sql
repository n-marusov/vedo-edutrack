-- create schema "identity_access"
CREATE SCHEMA IF NOT EXISTS identity_access;

-- create "roles"
CREATE TABLE IF NOT EXISTS identity_access.roles (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    name text NOT NULL,
    archetype text NOT NULL,
    description text NOT NULL DEFAULT '',
    PRIMARY KEY (id),
    CONSTRAINT roles_name_key UNIQUE (name)
);

-- create "role_permissions"
CREATE TABLE IF NOT EXISTS identity_access.role_permissions (
    role_id uuid NOT NULL,
    permission text NOT NULL,
    scope text NOT NULL DEFAULT 'own',
    PRIMARY KEY (role_id, permission, scope),
    CONSTRAINT role_permissions_role_id_fkey FOREIGN KEY (role_id) REFERENCES identity_access.roles (id) ON DELETE CASCADE
);

-- create "users"
CREATE TABLE IF NOT EXISTS identity_access.users (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    email text NOT NULL,
    display_name text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    CONSTRAINT users_email_key UNIQUE (email)
);

-- create "user_roles" (link users to roles)
CREATE TABLE IF NOT EXISTS identity_access.user_roles (
    user_id uuid NOT NULL,
    role_id uuid NOT NULL,
    PRIMARY KEY (user_id, role_id),
    CONSTRAINT user_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES identity_access.users (id) ON DELETE CASCADE,
    CONSTRAINT user_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES identity_access.roles (id) ON DELETE CASCADE
);
