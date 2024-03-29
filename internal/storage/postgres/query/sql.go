package query

// Users queries

const (
	InsertNewUser = `
INSERT INTO users
    (email, pass_hash) 
VALUES 
    ($1, $2)
RETURNING id
`

	SelectUserByEmail = `
SELECT 
	id, email, pass_hash, is_admin 
FROM users 
WHERE email = $1
`

	SelectIsAdminUserByID = `
SELECT 
	is_admin
FROM users 
WHERE 
    id=$1 AND is_admin IS TRUE
`
)

// App queries

const (
	SelectAppByID = `
SELECT 
	id, name, secret
FROM apps
WHERE id = $1
`
)
