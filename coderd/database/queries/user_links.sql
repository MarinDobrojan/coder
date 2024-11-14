-- name: GetUserLinkByLinkedID :one
SELECT
	user_links.*
FROM
	user_links
INNER JOIN
	users ON user_links.user_id = users.id
WHERE
	linked_id = $1
	AND
	deleted = false;

-- name: GetUserLinkByUserIDLoginType :one
SELECT
	*
FROM
	user_links
WHERE
	user_id = $1 AND login_type = $2;

-- name: GetUserLinksByUserID :many
SELECT * FROM user_links WHERE user_id = $1;

-- name: InsertUserLink :one
INSERT INTO
	user_links (
		user_id,
		login_type,
		linked_id,
		oauth_access_token,
		oauth_access_token_key_id,
		oauth_refresh_token,
		oauth_refresh_token_key_id,
		oauth_expiry,
		claims
	)
VALUES
	( $1, $2, $3, $4, $5, $6, $7, $8, $9 ) RETURNING *;

-- name: UpdateUserLinkedID :one
UPDATE
	user_links
SET
	linked_id = $1
WHERE
	user_id = $2 AND login_type = $3 RETURNING *;

-- name: UpdateUserLink :one
UPDATE
	user_links
SET
	oauth_access_token = $1,
	oauth_access_token_key_id = $2,
	oauth_refresh_token = $3,
	oauth_refresh_token_key_id = $4,
	oauth_expiry = $5,
	claims = $6
WHERE
	user_id = $7 AND login_type = $8 RETURNING *;


-- name: OIDCClaimFields :many
-- OIDCClaimFields returns a list of distinct keys in both the id_token_claims and user_info_claims fields.
-- This query is used to generate the list of available sync fields for idp sync settings.
SELECT
	DISTINCT jsonb_object_keys(claims->'id_token_claims')
FROM
	user_links
WHERE
    -- Only return rows where the top level key exists
	claims ? 'id_token_claims' AND
    -- 'null' is the default value for the id_token_claims field
	-- jsonb 'null' is not the same as SQL NULL. Strip these out.
	jsonb_typeof(claims->'id_token_claims') != 'null' AND
	login_type = 'oidc'
	AND CASE WHEN @organization_id :: uuid != '00000000-0000-0000-0000-000000000000'::uuid  THEN
		user_links.user_id = ANY(SELECT organization_members.user_id FROM organization_members WHERE organization_id = @organization_id)
		ELSE true
	END

-- Merge with user_info claims.
UNION

 -- This query is identical to the one above, except for 'user_info_claims'.
 -- There might be some way to do this more concisely at a cost of readability.
SELECT
	DISTINCT jsonb_object_keys(claims->'user_info_claims')
FROM
	user_links
WHERE
	claims ? 'user_info_claims' AND
	jsonb_typeof(claims->'user_info_claims') != 'null' AND
	login_type = 'oidc'
	AND CASE WHEN @organization_id :: uuid != '00000000-0000-0000-0000-000000000000'::uuid  THEN
		user_links.user_id = ANY(SELECT organization_members.user_id FROM organization_members WHERE organization_id = @organization_id)
		ELSE true
	END
;
