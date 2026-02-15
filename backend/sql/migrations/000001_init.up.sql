CREATE TABLE IF NOT EXISTS users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    nickname VARCHAR(100) NOT NULL DEFAULT '',
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_settings (
    user_id BIGINT PRIMARY KEY,
    llm_provider VARCHAR(50) NOT NULL DEFAULT '',
    llm_api_key TEXT NOT NULL,
    llm_base_url VARCHAR(500) NOT NULL DEFAULT '',
    llm_model VARCHAR(100) NOT NULL DEFAULT '',
    tts_provider VARCHAR(50) NOT NULL DEFAULT '',
    tts_api_key TEXT NOT NULL,
    tts_voice VARCHAR(100) NOT NULL DEFAULT '',
    tts_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    stt_provider VARCHAR(50) NOT NULL DEFAULT 'browser',
    stt_api_key TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_settings_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS interviews (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL DEFAULT '',
    position VARCHAR(255) NOT NULL DEFAULT '',
    status ENUM('pending', 'in_progress', 'completed') NOT NULL DEFAULT 'pending',
    language VARCHAR(10) NOT NULL DEFAULT 'zh-CN',
    llm_provider VARCHAR(50) NOT NULL DEFAULT '',
    llm_model VARCHAR(100) NOT NULL DEFAULT '',
    tts_provider VARCHAR(50) NOT NULL DEFAULT '',
    tts_voice VARCHAR(100) NOT NULL DEFAULT '',
    resume TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    CONSTRAINT fk_interviews_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS interview_messages (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    interview_id BIGINT NOT NULL,
    role ENUM('system', 'user', 'assistant') NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_interview_id (interview_id),
    CONSTRAINT fk_messages_interview FOREIGN KEY (interview_id) REFERENCES interviews(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS evaluations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    interview_id BIGINT NOT NULL UNIQUE,
    overall_score INT NOT NULL DEFAULT 0,
    summary TEXT NOT NULL,
    categories JSON,
    strengths TEXT NOT NULL DEFAULT '',
    weaknesses TEXT NOT NULL DEFAULT '',
    suggestions TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_evaluations_interview FOREIGN KEY (interview_id) REFERENCES interviews(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
