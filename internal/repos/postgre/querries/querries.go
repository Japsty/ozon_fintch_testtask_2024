package querries

const (
	// COMMENTS QUERRIES---------------------------------

	CreateComment = `
		INSERT INTO comments(id, content, author_id, post_id, parent_comment_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
	`
	GetCommentsByPostID = `
		SELECT id, content, author_id, post_id, parent_comment_id, created_at, updated_at
		FROM comments WHERE post_id = $1
	`
	//----------------------------------------------

	// POST QUERRIES---------------------------------

	GetAllPosts = `
		SELECT id, title, content, comments_allowed
		FROM posts;
	`
	CreatePost = `
		INSERT INTO posts (id, title, content, comments_allowed)
		VALUES ($1, $2, $3, $4);
	`
	GetPostByID = `
		SELECT id, title, content, comments_allowed
		FROM posts
		WHERE id = $1;
	`
	UpdatePostComments = `
		UPDATE posts
		SET comments_allowed = $2
		WHERE id = $1
		RETURNING id, title, content, comments_allowed;
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
