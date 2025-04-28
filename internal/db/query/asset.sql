-- name: CreateAsset :one
INSERT INTO assets (ticker, asset_name)
VALUES ($1, $2)
RETURNING id;

-- name: DeleteAsset :one
DELETE FROM assets
WHERE ticker = $1
RETURNING id;

-- name: UpdateAsset :one
UPDATE assets
SET asset_name = $1, is_tradable = $2
WHERE ticker = $3
RETURNING id;

-- name: GetAssetByTicker :one
SELECT * FROM assets
WHERE ticker = $1;

-- name: GetAllAssets :many
SELECT * FROM assets;
