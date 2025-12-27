-- name: GetTicketForUpdate :one
SELECT id, quantity FROM tickets 
WHERE id = $1 LIMIT 1 
FOR UPDATE;

-- name: UpdateTicket :exec
UPDATE tickets 
SET quantity = $2 
WHERE id = $1;

-- name: GetAllTicket :many
SELECT id, quantity FROM tickets;