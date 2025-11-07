START TRANSACTION;

INSERT INTO meeting_reservations (
  name,
  email,
  topic,
  message,
  start_at,
  end_at,
  duration_minutes,
  google_event_id,
  google_calendar_status,
  status,
  confirmation_sent_at,
  last_notification_sent_at,
  lookup_hash,
  cancellation_reason,
  created_at,
  updated_at
)
SELECT
  m.name,
  m.email,
  NULL,
  m.notes,
  m.meeting_at,
  DATE_ADD(m.meeting_at, INTERVAL COALESCE(m.duration_minutes, 0) MINUTE),
  COALESCE(m.duration_minutes, 0),
  m.calendar_event_id,
  CASE
    WHEN m.status = 'cancelled' THEN 'cancelled'
    ELSE 'confirmed'
  END,
  CASE
    WHEN m.status = 'cancelled' THEN 'cancelled'
    WHEN m.status = 'confirmed' THEN 'confirmed'
    ELSE 'pending'
  END,
  CASE
    WHEN m.status = 'confirmed' THEN m.updated_at
    ELSE NULL
  END,
  CASE
    WHEN m.status = 'confirmed' THEN m.updated_at
    ELSE NULL
  END,
  LOWER(SHA2(CONCAT_WS('-', COALESCE(m.email, ''), DATE_FORMAT(m.meeting_at, '%Y%m%d%H%i%s'), m.id), 256)),
  NULL,
  m.created_at,
  m.updated_at
FROM meetings m
WHERE NOT EXISTS (SELECT 1 FROM meeting_reservations LIMIT 1);

COMMIT;
