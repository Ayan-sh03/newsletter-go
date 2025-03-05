-- User Management Queries
-- Create a new user
-- name: CreateUser :one
INSERT INTO users (username, password, email)
VALUES ($1, $2, $3)
RETURNING *;

-- Get user by ID
-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND status = 'active';

-- Get user by username
-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 AND status = 'active';

-- Get user by email
-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND status = 'active';

-- Update user
-- name: UpdateUser :one
UPDATE users
SET username = COALESCE($2, username),
    email = COALESCE($3, email),
    updated_at = current_timestamp
WHERE id = $1
RETURNING *;

-- Update user password
-- name: UpdateUserPassword :exec
UPDATE users
SET password = $2,
    updated_at = current_timestamp
WHERE id = $1;

-- Update user last login
-- name: UpdateUserLastLogin :exec
UPDATE users
SET last_login = current_timestamp
WHERE id = $1;

-- Deactivate user
-- name: DeactivateUser :exec
UPDATE users
SET status = 'inactive',
    updated_at = current_timestamp
WHERE id = $1;

-- Newsletter Management Queries
-- Create a newsletter
-- name: CreateNewsletter :one
INSERT INTO newsletters (title, author, description)
VALUES ($1, $2, $3)
RETURNING *;

-- Get newsletter by ID
-- name: GetNewsletterByID :one
SELECT n.*, u.username as author_name
FROM newsletters n
JOIN users u ON n.author = u.id
WHERE n.id = $1;

-- Get newsletters by author
-- name: GetNewslettersByAuthor :many
SELECT n.*, u.username as author_name
FROM newsletters n
JOIN users u ON n.author = u.id
WHERE n.author = $1
ORDER BY n.created_at DESC;

-- Update newsletter
-- name: UpdateNewsletter :one
UPDATE newsletters
SET title = COALESCE($2, title),
    description = COALESCE($3, description),
    updated_at = current_timestamp
WHERE id = $1
RETURNING *;

-- Delete newsletter
-- name: DeleteNewsletter :exec
DELETE FROM newsletters
WHERE id = $1;

-- Get all newsletters with pagination
-- name: ListNewsletters :many
SELECT n.*, u.username as author_name
FROM newsletters n
JOIN users u ON n.author = u.id
ORDER BY n.created_at DESC
LIMIT $1 OFFSET $2;

-- Subscription Management Queries
-- Subscribe user to newsletter
-- name: SubscribeToNewsletter :one
INSERT INTO subscribers (user_id, newsletter_id)
VALUES ($1, $2)
RETURNING *;

-- Unsubscribe user from newsletter
-- name: UnsubscribeFromNewsletter :exec
DELETE FROM subscribers
WHERE user_id = $1 AND newsletter_id = $2;

-- Get all subscriptions for a user
-- name: GetUserSubscriptions :many
SELECT s.*, n.title as newsletter_title
FROM subscribers s
JOIN newsletters n ON s.newsletter_id = n.id
WHERE s.user_id = $1
ORDER BY s.created_at DESC;

-- Get all subscribers for a newsletter
-- name: GetNewsletterSubscribers :many
SELECT s.*, u.username, u.email
FROM subscribers s
JOIN users u ON s.user_id = u.id
WHERE s.newsletter_id = $1
ORDER BY s.created_at DESC;

-- Check if user is subscribed to newsletter
-- name: IsUserSubscribed :one
SELECT EXISTS (
    SELECT 1 FROM subscribers
    WHERE user_id = $1 AND newsletter_id = $2
) AS is_subscribed;

-- Get subscriber count for a newsletter
-- name: GetNewsletterSubscriberCount :one
SELECT COUNT(*) FROM subscribers
WHERE newsletter_id = $1;

-- Letter Management Queries
-- Create a letter
-- name: CreateLetter :one
INSERT INTO letters (newsletter_id, title, content, status)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- Get letter by ID
-- name: GetLetterByID :one
SELECT l.*, n.title as newsletter_title
FROM letters l
JOIN newsletters n ON l.newsletter_id = n.id
WHERE l.id = $1;

-- Get all letters for a newsletter
-- name: GetNewsletterLetters :many
SELECT *
FROM letters
WHERE newsletter_id = $1
ORDER BY created_at DESC;

-- Update letter
-- name: UpdateLetter :one
UPDATE letters
SET title = COALESCE($2, title),
    content = COALESCE($3, content),
    status = COALESCE($4, status),
    updated_at = current_timestamp
WHERE id = $1
RETURNING *;

-- Publish letter (change status from draft to published)
-- name: PublishLetter :one
UPDATE letters
SET status = 'published',
    updated_at = current_timestamp
WHERE id = $1
RETURNING *;

-- Delete letter
-- name: DeleteLetter :exec
DELETE FROM letters
WHERE id = $1;

-- View Management Queries
-- Record a view
-- name: RecordView :one
INSERT INTO views (letter_id, user_id, ip_address)
VALUES ($1, $2, $3)
RETURNING *;

-- Get view count for a letter
-- name: GetLetterViewCount :one
SELECT COUNT(*) FROM views
WHERE letter_id = $1;

-- Get unique viewer count for a letter
-- name: GetLetterUniqueViewerCount :one
SELECT COUNT(DISTINCT user_id) FROM views
WHERE letter_id = $1;

-- Get views by user
-- name: GetViewsByUser :many
SELECT v.*, l.title as letter_title
FROM views v
JOIN letters l ON v.letter_id = l.id
WHERE v.user_id = $1
ORDER BY v.created_at DESC;

-- Analytics Queries
-- Get most viewed letters
-- name: GetMostViewedLetters :many
SELECT l.*, COUNT(v.id) as view_count
FROM letters l
JOIN views v ON l.id = v.letter_id
GROUP BY l.id
ORDER BY view_count DESC
LIMIT $1;

-- Get most active subscribers
-- name: GetMostActiveSubscribers :many
SELECT u.id, u.username, COUNT(v.id) as view_count
FROM users u
JOIN views v ON u.id = v.user_id
GROUP BY u.id
ORDER BY view_count DESC
LIMIT $1;

-- Get newsletter engagement stats
-- name: GetNewsletterEngagementStats :one
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
GROUP BY n.id;

-- Get recent activity for a user
-- name: GetUserRecentActivity :many
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
LIMIT $2;