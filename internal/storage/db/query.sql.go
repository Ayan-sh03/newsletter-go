// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package db

import (
	"context"
	"time"

	"github.com/sqlc-dev/pqtype"
)

const createLetter = `-- name: CreateLetter :one
INSERT INTO letters (newsletter_id, title, content, status)
VALUES ($1, $2, $3, $4)
RETURNING id, newsletter_id, title, content, status, created_at, updated_at
`

type CreateLetterParams struct {
	NewsletterID int64
	Title        string
	Content      string
	Status       string
}

// Letter Management Queries
// Create a letter
func (q *Queries) CreateLetter(ctx context.Context, arg CreateLetterParams) (Letter, error) {
	row := q.db.QueryRowContext(ctx, createLetter,
		arg.NewsletterID,
		arg.Title,
		arg.Content,
		arg.Status,
	)
	var i Letter
	err := row.Scan(
		&i.ID,
		&i.NewsletterID,
		&i.Title,
		&i.Content,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createNewsletter = `-- name: CreateNewsletter :one
INSERT INTO newsletters (title, author, description)
VALUES ($1, $2, $3)
RETURNING title, author, description
`

type CreateNewsletterParams struct {
	Title       string
	Author      int64
	Description string
}

// Newsletter Management Queries
// Create a newsletter
func (q *Queries) CreateNewsletter(ctx context.Context, arg CreateNewsletterParams) (Newsletter, error) {
	row := q.db.QueryRowContext(ctx, createNewsletter, arg.Title, arg.Author, arg.Description)
	var i Newsletter
	err := row.Scan(
		&i.Title,
		&i.Author,
		&i.Description,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, password, email)
VALUES ($1, $2, $3)
RETURNING id, username, password, email, status, created_at, updated_at, last_login
`

type CreateUserParams struct {
	Username string
	Password string
	Email    string
}

// User Management Queries
// Create a new user
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.Password, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastLogin,
	)
	return i, err
}

const deactivateUser = `-- name: DeactivateUser :exec
UPDATE users
SET status = 'inactive',
    updated_at = current_timestamp
WHERE id = $1
`

// Deactivate user
func (q *Queries) DeactivateUser(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deactivateUser, id)
	return err
}

const deleteLetter = `-- name: DeleteLetter :exec
DELETE FROM letters
WHERE id = $1
`

// Delete letter
func (q *Queries) DeleteLetter(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteLetter, id)
	return err
}

const deleteNewsletter = `-- name: DeleteNewsletter :exec
DELETE FROM newsletters
WHERE id = $1
`

// Delete newsletter
func (q *Queries) DeleteNewsletter(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteNewsletter, id)
	return err
}

const getLetterByID = `-- name: GetLetterByID :one
SELECT l.id, l.newsletter_id, l.title, l.content, l.status, l.created_at, l.updated_at, n.title as newsletter_title
FROM letters l
JOIN newsletters n ON l.newsletter_id = n.id
WHERE l.id = $1
`

type GetLetterByIDRow struct {
	ID              int64
	NewsletterID    int64
	Title           string
	Content         string
	Status          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	NewsletterTitle string
}

// Get letter by ID
func (q *Queries) GetLetterByID(ctx context.Context, id int64) (GetLetterByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getLetterByID, id)
	var i GetLetterByIDRow
	err := row.Scan(
		&i.ID,
		&i.NewsletterID,
		&i.Title,
		&i.Content,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.NewsletterTitle,
	)
	return i, err
}

const getLetterUniqueViewerCount = `-- name: GetLetterUniqueViewerCount :one
SELECT COUNT(DISTINCT user_id) FROM views
WHERE letter_id = $1
`

// Get unique viewer count for a letter
func (q *Queries) GetLetterUniqueViewerCount(ctx context.Context, letterID int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getLetterUniqueViewerCount, letterID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getLetterViewCount = `-- name: GetLetterViewCount :one
SELECT COUNT(*) FROM views
WHERE letter_id = $1
`

// Get view count for a letter
func (q *Queries) GetLetterViewCount(ctx context.Context, letterID int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getLetterViewCount, letterID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getMostActiveSubscribers = `-- name: GetMostActiveSubscribers :many
SELECT u.id, u.username, COUNT(v.id) as view_count
FROM users u
JOIN views v ON u.id = v.user_id
GROUP BY u.id
ORDER BY view_count DESC
LIMIT $1
`

type GetMostActiveSubscribersRow struct {
	ID        int64
	Username  string
	ViewCount int64
}

// Get most active subscribers
func (q *Queries) GetMostActiveSubscribers(ctx context.Context, limit int32) ([]GetMostActiveSubscribersRow, error) {
	rows, err := q.db.QueryContext(ctx, getMostActiveSubscribers, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMostActiveSubscribersRow
	for rows.Next() {
		var i GetMostActiveSubscribersRow
		if err := rows.Scan(&i.ID, &i.Username, &i.ViewCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMostViewedLetters = `-- name: GetMostViewedLetters :many
SELECT l.id, l.newsletter_id, l.title, l.content, l.status, l.created_at, l.updated_at, COUNT(v.id) as view_count
FROM letters l
JOIN views v ON l.id = v.letter_id
GROUP BY l.id
ORDER BY view_count DESC
LIMIT $1
`

type GetMostViewedLettersRow struct {
	ID           int64
	NewsletterID int64
	Title        string
	Content      string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ViewCount    int64
}

// Analytics Queries
// Get most viewed letters
func (q *Queries) GetMostViewedLetters(ctx context.Context, limit int32) ([]GetMostViewedLettersRow, error) {
	rows, err := q.db.QueryContext(ctx, getMostViewedLetters, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMostViewedLettersRow
	for rows.Next() {
		var i GetMostViewedLettersRow
		if err := rows.Scan(
			&i.ID,
			&i.NewsletterID,
			&i.Title,
			&i.Content,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ViewCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNewsletterByID = `-- name: GetNewsletterByID :one
SELECT n.id, n.title, n.author, n.description, n.status, n.created_at, n.updated_at, u.username as author_name
FROM newsletters n
JOIN users u ON n.author = u.id
WHERE n.id = $1
`

type GetNewsletterByIDRow struct {
	ID          int64
	Title       string
	Author      int64
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AuthorName  string
}

// Get newsletter by ID
func (q *Queries) GetNewsletterByID(ctx context.Context, id int64) (GetNewsletterByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getNewsletterByID, id)
	var i GetNewsletterByIDRow
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Author,
		&i.Description,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.AuthorName,
	)
	return i, err
}

const getNewsletterEngagementStats = `-- name: GetNewsletterEngagementStats :one
SELECT 
    n.id,
    n.title,
    COUNT(DISTINCT s.user_id) as subscriber_count,
    COUNT(DISTINCT l.id) as letter_count,
    COUNT(DISTINCT v.id) as total_views
FROM newsletters n
LEFT JOIN subscribers s ON n.id = s.newsletter_id
LEFT JOIN letters l ON n.id = l.newsletter_id
LEFT JOIN views v ON l.id = v.letter_id
WHERE n.id = $1
GROUP BY n.id
`

type GetNewsletterEngagementStatsRow struct {
	ID              int64
	Title           string
	SubscriberCount int64
	LetterCount     int64
	TotalViews      int64
}

// Get newsletter engagement stats
func (q *Queries) GetNewsletterEngagementStats(ctx context.Context, id int64) (GetNewsletterEngagementStatsRow, error) {
	row := q.db.QueryRowContext(ctx, getNewsletterEngagementStats, id)
	var i GetNewsletterEngagementStatsRow
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.SubscriberCount,
		&i.LetterCount,
		&i.TotalViews,
	)
	return i, err
}

const getNewsletterLetters = `-- name: GetNewsletterLetters :many
SELECT id, newsletter_id, title, content, status, created_at, updated_at
FROM letters
WHERE newsletter_id = $1
ORDER BY created_at DESC
`

// Get all letters for a newsletter
func (q *Queries) GetNewsletterLetters(ctx context.Context, newsletterID int64) ([]Letter, error) {
	rows, err := q.db.QueryContext(ctx, getNewsletterLetters, newsletterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Letter
	for rows.Next() {
		var i Letter
		if err := rows.Scan(
			&i.ID,
			&i.NewsletterID,
			&i.Title,
			&i.Content,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNewsletterSubscriberCount = `-- name: GetNewsletterSubscriberCount :one
SELECT COUNT(*) FROM subscribers
WHERE newsletter_id = $1
`

// Get subscriber count for a newsletter
func (q *Queries) GetNewsletterSubscriberCount(ctx context.Context, newsletterID int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getNewsletterSubscriberCount, newsletterID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getNewsletterSubscribers = `-- name: GetNewsletterSubscribers :many
SELECT s.id, s.user_id, s.newsletter_id, s.created_at, s.updated_at, u.username, u.email
FROM subscribers s
JOIN users u ON s.user_id = u.id
WHERE s.newsletter_id = $1
ORDER BY s.created_at DESC
`

type GetNewsletterSubscribersRow struct {
	ID           int64
	UserID       int64
	NewsletterID int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Username     string
	Email        string
}

// Get all subscribers for a newsletter
func (q *Queries) GetNewsletterSubscribers(ctx context.Context, newsletterID int64) ([]GetNewsletterSubscribersRow, error) {
	rows, err := q.db.QueryContext(ctx, getNewsletterSubscribers, newsletterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetNewsletterSubscribersRow
	for rows.Next() {
		var i GetNewsletterSubscribersRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.NewsletterID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Username,
			&i.Email,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNewslettersByAuthor = `-- name: GetNewslettersByAuthor :many
SELECT n.id, n.title, n.author, n.description, n.status, n.created_at, n.updated_at, u.username as author_name
FROM newsletters n
JOIN users u ON n.author = u.id
WHERE n.author = $1
ORDER BY n.created_at DESC
`

type GetNewslettersByAuthorRow struct {
	ID          int64
	Title       string
	Author      int64
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AuthorName  string
}

// Get newsletters by author
func (q *Queries) GetNewslettersByAuthor(ctx context.Context, author int64) ([]GetNewslettersByAuthorRow, error) {
	rows, err := q.db.QueryContext(ctx, getNewslettersByAuthor, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetNewslettersByAuthorRow
	for rows.Next() {
		var i GetNewslettersByAuthorRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Author,
			&i.Description,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.AuthorName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, username, password, email, status, created_at, updated_at, last_login FROM users
WHERE email = $1 AND status = 'active'
`

// Get user by email
func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastLogin,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, username, password, email, status, created_at, updated_at, last_login FROM users
WHERE id = $1 AND status = 'active'
`

// Get user by ID
func (q *Queries) GetUserByID(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastLogin,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, password, email, status, created_at, updated_at, last_login FROM users
WHERE username = $1 AND status = 'active'
`

// Get user by username
func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastLogin,
	)
	return i, err
}

const getUserRecentActivity = `-- name: GetUserRecentActivity :many
SELECT 
    'view' as activity_type,
    v.created_at as activity_time,
    l.title as target_title,
    n.id as newsletter_id,
    n.title as newsletter_title
FROM views v
JOIN letters l ON v.letter_id = l.id
JOIN newsletters n ON l.newsletter_id = n.id
WHERE v.user_id = $1
UNION ALL
SELECT 
    'subscribe' as activity_type,
    s.created_at as activity_time,
    n.title as target_title,
    n.id as newsletter_id,
    n.title as newsletter_title
FROM subscribers s
JOIN newsletters n ON s.newsletter_id = n.id
WHERE s.user_id = $1
ORDER BY activity_time DESC
LIMIT $2
`

type GetUserRecentActivityParams struct {
	UserID int64
	Limit  int32
}

type GetUserRecentActivityRow struct {
	ActivityType    string
	ActivityTime    time.Time
	TargetTitle     string
	NewsletterID    int64
	NewsletterTitle string
}

// Get recent activity for a user
func (q *Queries) GetUserRecentActivity(ctx context.Context, arg GetUserRecentActivityParams) ([]GetUserRecentActivityRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserRecentActivity, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserRecentActivityRow
	for rows.Next() {
		var i GetUserRecentActivityRow
		if err := rows.Scan(
			&i.ActivityType,
			&i.ActivityTime,
			&i.TargetTitle,
			&i.NewsletterID,
			&i.NewsletterTitle,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserSubscriptions = `-- name: GetUserSubscriptions :many
SELECT s.id, s.user_id, s.newsletter_id, s.created_at, s.updated_at, n.title as newsletter_title
FROM subscribers s
JOIN newsletters n ON s.newsletter_id = n.id
WHERE s.user_id = $1
ORDER BY s.created_at DESC
`

type GetUserSubscriptionsRow struct {
	ID              int64
	UserID          int64
	NewsletterID    int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	NewsletterTitle string
}

// Get all subscriptions for a user
func (q *Queries) GetUserSubscriptions(ctx context.Context, userID int64) ([]GetUserSubscriptionsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserSubscriptions, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserSubscriptionsRow
	for rows.Next() {
		var i GetUserSubscriptionsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.NewsletterID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.NewsletterTitle,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getViewsByUser = `-- name: GetViewsByUser :many
SELECT v.id, v.letter_id, v.user_id, v.ip_address, v.created_at, v.updated_at, l.title as letter_title
FROM views v
JOIN letters l ON v.letter_id = l.id
WHERE v.user_id = $1
ORDER BY v.created_at DESC
`

type GetViewsByUserRow struct {
	ID          int64
	LetterID    int64
	UserID      int64
	IpAddress   pqtype.Inet
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LetterTitle string
}

// Get views by user
func (q *Queries) GetViewsByUser(ctx context.Context, userID int64) ([]GetViewsByUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getViewsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetViewsByUserRow
	for rows.Next() {
		var i GetViewsByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.LetterID,
			&i.UserID,
			&i.IpAddress,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.LetterTitle,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const isUserSubscribed = `-- name: IsUserSubscribed :one
SELECT EXISTS (
    SELECT 1 FROM subscribers
    WHERE user_id = $1 AND newsletter_id = $2
) AS is_subscribed
`

type IsUserSubscribedParams struct {
	UserID       int64
	NewsletterID int64
}

// Check if user is subscribed to newsletter
func (q *Queries) IsUserSubscribed(ctx context.Context, arg IsUserSubscribedParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, isUserSubscribed, arg.UserID, arg.NewsletterID)
	var is_subscribed bool
	err := row.Scan(&is_subscribed)
	return is_subscribed, err
}

const listNewsletters = `-- name: ListNewsletters :many
SELECT n.id, n.title, n.author, n.description, n.status, n.created_at, n.updated_at, u.username as author_name
FROM newsletters n
JOIN users u ON n.author = u.id
ORDER BY n.created_at DESC
LIMIT $1 OFFSET $2
`

type ListNewslettersParams struct {
	Limit  int32
	Offset int32
}

type ListNewslettersRow struct {
	ID          int64
	Title       string
	Author      int64
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AuthorName  string
}

// Get all newsletters with pagination
func (q *Queries) ListNewsletters(ctx context.Context, arg ListNewslettersParams) ([]ListNewslettersRow, error) {
	rows, err := q.db.QueryContext(ctx, listNewsletters, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListNewslettersRow
	for rows.Next() {
		var i ListNewslettersRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Author,
			&i.Description,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.AuthorName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const publishLetter = `-- name: PublishLetter :one
UPDATE letters
SET status = 'published',
    updated_at = current_timestamp
WHERE id = $1
RETURNING id, newsletter_id, title, content, status, created_at, updated_at
`

// Publish letter (change status from draft to published)
func (q *Queries) PublishLetter(ctx context.Context, id int64) (Letter, error) {
	row := q.db.QueryRowContext(ctx, publishLetter, id)
	var i Letter
	err := row.Scan(
		&i.ID,
		&i.NewsletterID,
		&i.Title,
		&i.Content,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const recordView = `-- name: RecordView :one
INSERT INTO views (letter_id, user_id, ip_address)
VALUES ($1, $2, $3)
RETURNING id, letter_id, user_id, ip_address, created_at, updated_at
`

type RecordViewParams struct {
	LetterID  int64
	UserID    int64
	IpAddress pqtype.Inet
}

// View Management Queries
// Record a view
func (q *Queries) RecordView(ctx context.Context, arg RecordViewParams) (View, error) {
	row := q.db.QueryRowContext(ctx, recordView, arg.LetterID, arg.UserID, arg.IpAddress)
	var i View
	err := row.Scan(
		&i.ID,
		&i.LetterID,
		&i.UserID,
		&i.IpAddress,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const subscribeToNewsletter = `-- name: SubscribeToNewsletter :one
INSERT INTO subscribers (user_id, newsletter_id)
VALUES ($1, $2)
RETURNING id, user_id, newsletter_id, created_at, updated_at
`

type SubscribeToNewsletterParams struct {
	UserID       int64
	NewsletterID int64
}

// Subscription Management Queries
// Subscribe user to newsletter
func (q *Queries) SubscribeToNewsletter(ctx context.Context, arg SubscribeToNewsletterParams) (Subscriber, error) {
	row := q.db.QueryRowContext(ctx, subscribeToNewsletter, arg.UserID, arg.NewsletterID)
	var i Subscriber
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.NewsletterID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const unsubscribeFromNewsletter = `-- name: UnsubscribeFromNewsletter :exec
DELETE FROM subscribers
WHERE user_id = $1 AND newsletter_id = $2
`

type UnsubscribeFromNewsletterParams struct {
	UserID       int64
	NewsletterID int64
}

// Unsubscribe user from newsletter
func (q *Queries) UnsubscribeFromNewsletter(ctx context.Context, arg UnsubscribeFromNewsletterParams) error {
	_, err := q.db.ExecContext(ctx, unsubscribeFromNewsletter, arg.UserID, arg.NewsletterID)
	return err
}

const updateLetter = `-- name: UpdateLetter :one
UPDATE letters
SET title = COALESCE($2, title),
    content = COALESCE($3, content),
    status = COALESCE($4, status),
    updated_at = current_timestamp
WHERE id = $1
RETURNING id, newsletter_id, title, content, status, created_at, updated_at
`

type UpdateLetterParams struct {
	ID      int64
	Title   string
	Content string
	Status  string
}

// Update letter
func (q *Queries) UpdateLetter(ctx context.Context, arg UpdateLetterParams) (Letter, error) {
	row := q.db.QueryRowContext(ctx, updateLetter,
		arg.ID,
		arg.Title,
		arg.Content,
		arg.Status,
	)
	var i Letter
	err := row.Scan(
		&i.ID,
		&i.NewsletterID,
		&i.Title,
		&i.Content,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateNewsletter = `-- name: UpdateNewsletter :one
UPDATE newsletters
SET title = COALESCE($2, title),
    description = COALESCE($3, description),
    updated_at = current_timestamp
WHERE id = $1
RETURNING id, title, author, description, status, created_at, updated_at
`

type UpdateNewsletterParams struct {
	ID          int64
	Title       string
	Description string
}

// Update newsletter
func (q *Queries) UpdateNewsletter(ctx context.Context, arg UpdateNewsletterParams) (Newsletter, error) {
	row := q.db.QueryRowContext(ctx, updateNewsletter, arg.ID, arg.Title, arg.Description)
	var i Newsletter
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Author,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET username = COALESCE($2, username),
    email = COALESCE($3, email),
    updated_at = current_timestamp
WHERE id = $1
RETURNING id, username, password, email, status, created_at, updated_at, last_login
`

type UpdateUserParams struct {
	ID       int64
	Username string
	Email    string
}

// Update user
func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser, arg.ID, arg.Username, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastLogin,
	)
	return i, err
}

const updateUserLastLogin = `-- name: UpdateUserLastLogin :exec
UPDATE users
SET last_login = current_timestamp
WHERE id = $1
`

// Update user last login
func (q *Queries) UpdateUserLastLogin(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, updateUserLastLogin, id)
	return err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2,
    updated_at = current_timestamp
WHERE id = $1
`

type UpdateUserPasswordParams struct {
	ID       int64
	Password string
}

// Update user password
func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.ExecContext(ctx, updateUserPassword, arg.ID, arg.Password)
	return err
}
