CREATE TABLE IF NOT EXISTS meetings (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(320) NOT NULL,
    meeting_at DATETIME NOT NULL,
    duration_minutes INT NOT NULL DEFAULT 30,
    meet_url VARCHAR(512) NOT NULL DEFAULT '',
    calendar_event_id VARCHAR(512) NOT NULL DEFAULT '',
    status ENUM('pending', 'confirmed', 'cancelled') NOT NULL DEFAULT 'pending',
    notes TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_meetings_meeting_at (meeting_at),
    INDEX idx_meetings_status_meeting_at (status, meeting_at)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS schedule_blackouts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    reason VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_blackouts_start (start_time),
    INDEX idx_blackouts_end (end_time)
) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

INSERT INTO schedule_blackouts (start_time, end_time, reason)
VALUES
    ('2025-10-24 00:00:00', '2025-10-24 02:00:00', 'System maintenance window')
ON DUPLICATE KEY UPDATE
    start_time = VALUES(start_time),
    end_time = VALUES(end_time),
    reason = VALUES(reason);
