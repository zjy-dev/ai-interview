-- name: CreateInterview :execlastid
INSERT INTO interviews (user_id, title, position, status, language,
    llm_provider, llm_model, tts_provider, tts_voice, resume)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetInterviewByID :one
SELECT id, user_id, title, position, status, language,
    llm_provider, llm_model, tts_provider, tts_voice, resume, created_at, updated_at
FROM interviews WHERE id = ?;

-- name: ListInterviewsByUserID :many
SELECT id, user_id, title, position, status, language, created_at, updated_at
FROM interviews WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- name: CountInterviewsByUserID :one
SELECT COUNT(*) FROM interviews WHERE user_id = ?;

-- name: UpdateInterviewStatus :exec
UPDATE interviews SET status = ? WHERE id = ?;

-- name: CreateMessage :execlastid
INSERT INTO interview_messages (interview_id, role, content) VALUES (?, ?, ?);

-- name: ListMessagesByInterviewID :many
SELECT id, interview_id, role, content, created_at
FROM interview_messages WHERE interview_id = ? ORDER BY id ASC;

-- name: CreateEvaluation :execlastid
INSERT INTO evaluations (interview_id, overall_score, summary, categories, strengths, weaknesses, suggestions)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetEvaluationByInterviewID :one
SELECT id, interview_id, overall_score, summary, categories, strengths, weaknesses, suggestions, created_at
FROM evaluations WHERE interview_id = ?;
