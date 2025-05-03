-- name: CreateVacancy :one
INSERT INTO vacancies (id, created_at, updated_at, title, company_name, url)
VALUES (
	gen_random_uuid(),
	NOW(),
	NOW(),
	$1,
	$2,
	$3
)
RETURNING *;

-- name: GetAllVacancies :many
SELECT * FROM vacancies;

