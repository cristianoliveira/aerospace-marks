-- name: AddMark :exec
INSERT INTO marks (window_id, mark) VALUES (?, ?);

-- name: GetAllMarks :many
SELECT window_id, mark FROM marks;

-- name: GetMarksByWindowID :many
SELECT window_id, mark
FROM marks
WHERE window_id = ?;

-- name: GetWindowByMark :one
SELECT window_id, mark FROM marks WHERE mark = ?;

-- name: DeleteAllMarks :execresult
DELETE FROM marks;

-- name: DeleteByMark :execresult
DELETE FROM marks WHERE mark = ?;

-- name: DeleteByWindow :execresult
DELETE FROM marks WHERE window_id = ?;

-- name: DeleteMarksByWindowIDOrMark :execresult
DELETE FROM marks WHERE window_id = ? OR mark = ?;