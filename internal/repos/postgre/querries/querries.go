package querries

const (
	// SESSIONS QUERRIES---------------------------------
	CreateSession = `
		INSERT INTO sessions (user_id,username,exp,token) 
		VALUES (?, ?, ?, ?)
	`

	GetSessionByToken = `
		SELECT id, user_id, username, exp 
		FROM sessions 
		WHERE token = ?
    `

	DeleteSession = `
		DELETE FROM sessions 
		WHERE token = ?
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

	// COMMENTS QUERRIES---------------------------------

	CreateComment = `
		INSERT INTO comments(id, content, author_id, post_id, parent_id, created_at)
		VALUES ($1, $2, $3, $4, COALESCE(NULLIF($5, '')), CURRENT_TIMESTAMP)
		RETURNING id,content,author_id,post_id,parent_id,created_at;
	`

	GetAllCommentsByPostID = `
		SELECT id, content, author_id, post_id, parent_id, created_at
		FROM comments
		WHERE post_id = $1;
	`
	//----------------------------------------------

	// POST QUERRIES---------------------------------

	GetAllPosts = `
		SELECT id, title, content, user_id, comments_allowed, created_at
		FROM posts;
	`

	CreatePost = `
		INSERT INTO posts (id, title, content, user_id, comments_allowed,created_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
		RETURNING id, title, content, user_id, comments_allowed,created_at;
	`
	GetPostByID = `
		SELECT id, title, content, user_id, comments_allowed, created_at
		FROM posts
		WHERE id = $1;
	`
	UpdatePostCommentsStatus = `
		UPDATE posts
		SET comments_allowed = $2
		WHERE id = $1
		RETURNING id, title, content, user_id, comments_allowed,created_at;
	`

	//----------------------------------------------
)
