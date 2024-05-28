package querries

const (
	// COMMENTS QUERRIES---------------------------------

	CheckUserExists = `
		SELECT COUNT(*) > 0
		FROM users
		WHERE username = ?;
	`
	GetUserByLogin = `
		SELECT username,password_hash,user_id
		FROM users
		WHERE username = ?
	`

	CreateUser = `
		INSERT INTO users (user_id, username, password_hash)
		VALUES (?, ?, ?);
	`

	GetUserByID = `
		SELECT user_id, username, password_hash
		FROM users
		WHERE user_id = ?
    `
	//----------------------------------------------

	// POST QUERRIES---------------------------------

	CheckUserExists = `
		SELECT COUNT(*) > 0
		FROM users
		WHERE username = ?;
	`
	GetUserByLogin = `
		SELECT username,password_hash,user_id
		FROM users
		WHERE username = ?
	`

	CreateUser = `
		INSERT INTO users (user_id, username, password_hash)
		VALUES (?, ?, ?);
	`

	GetUserByID = `
		SELECT user_id, username, password_hash
		FROM users
		WHERE user_id = ?
    `
	//----------------------------------------------
	// USER QUERRIES---------------------------------

	CheckUserExists = `
		SELECT COUNT(*) > 0
		FROM users
		WHERE username = ?;
	`
	GetUserByLogin = `
		SELECT username,password_hash,user_id
		FROM users
		WHERE username = ?
	`

	CreateUser = `
		INSERT INTO users (user_id, username, password_hash)
		VALUES (?, ?, ?);
	`

	GetUserByID = `
		SELECT user_id, username, password_hash
		FROM users
		WHERE user_id = ?
    `
	//----------------------------------------------
)
