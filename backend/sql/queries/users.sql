-- name: CreateUser :execlastid
INSERT INTO users (email, nickname, password_hash) VALUES (?, ?, ?);

-- name: GetUserByID :one
SELECT id, email, nickname, password_hash, created_at, updated_at
FROM users WHERE id = ?;

-- name: GetUserByEmail :one
SELECT id, email, nickname, password_hash, created_at, updated_at
FROM users WHERE email = ?;

-- name: UpsertUserSettings :exec
INSERT INTO user_settings (user_id, llm_provider, llm_api_key, llm_base_url, llm_model,
    tts_provider, tts_api_key, tts_voice, tts_enabled, stt_provider, stt_api_key)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
    llm_provider = VALUES(llm_provider),
    llm_api_key = VALUES(llm_api_key),
    llm_base_url = VALUES(llm_base_url),
    llm_model = VALUES(llm_model),
    tts_provider = VALUES(tts_provider),
    tts_api_key = VALUES(tts_api_key),
    tts_voice = VALUES(tts_voice),
    tts_enabled = VALUES(tts_enabled),
    stt_provider = VALUES(stt_provider),
    stt_api_key = VALUES(stt_api_key);

-- name: GetUserSettings :one
SELECT user_id, llm_provider, llm_api_key, llm_base_url, llm_model,
    tts_provider, tts_api_key, tts_voice, tts_enabled, stt_provider, stt_api_key
FROM user_settings WHERE user_id = ?;
