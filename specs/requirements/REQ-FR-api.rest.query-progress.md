# REQ-FR-api.rest.query-progress

**Приоритет:** P1
**Ключевая функция:** F6.1
**Источник:** UC-api.rest.query-progress

**Описание:** Сервис должен предоставлять REST API endpoint `GET /progress/{learner_id}` для запроса прогресса конкретного обучающегося: освоенные модули, статус активной траектории, процент выполнения плана и временные метки обновления. Endpoint доступен только аутентифицированным и авторизованным субъектам (learner, parent/HR, methodologist, LMS-платформа с подпиской). Ответ должен содержать все данные, необходимые для отображения прогресса без дополнительных запросов.

**Критерии приёмки:**
- `GET /progress/{learner_id}` возвращает `200 OK` с JSON-схемой `LearnerProgress` (поля `learner_id`, `plan_id`, `status` (`active|completed|deviated`), `mastered_modules[]`, `total_progress_percent` (0–100), `last_updated` в ISO 8601); все обязательные поля присутствуют и валидны по схеме.
- Запрос без валидного `Authorization` header (Bearer token, выданный через SSO/Keycloak) возвращает `401 Unauthorized`; запрос от субъекта без права доступа к данному `learner_id` возвращает `403 Forbidden`.
- Запрос для несуществующего `learner_id` возвращает `404 Not Found` с идентификатором ошибки в теле ответа; `learner_id` в невалидном формате (UUID) возвращает `400 Bad Request`.
- Время ответа p95 не превышает 500 мс на эталонном окружении при отсутствии кэша (cross-ref REQ-NFR-api.performance.progress-query-latency).
