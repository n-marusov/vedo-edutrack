<a id="us-api.rest.query-progress"></a>
# US-api.rest.query-progress: REST-эндпоинт запроса прогресса

@US-api.rest.query-progress @UC-api.rest.query-progress @P1 @api
Feature: US-api.rest.query-progress Запрос прогресса Learner через REST API

  Background:
    Given внешняя система аутентифицирована по API-ключу
    Given API предоставляет эндпоинт GET /api/v1/learners/{id}/progress

  Scenario: Запрос прогресса Learner
    Given Learner имеет активную Route
    When внешняя система запросила прогресс Learner
    Then API возвращает статус 200
    And ответ содержит освоенные Checkpoint и текущий уровень Mastery

  Scenario: Learner не найден
    Given Learner с указанным id не существует
    When внешняя система запросила прогресс
    Then API возвращает ошибку "Not Found: learner does not exist"

  Scenario: Актуальность прогресса после события
    Given система получила событие domain "module.mastered"
    When событие обработано
    Then прогресс Learner обновлён
    And последующие запросы прогресса возвращают актуальные данные
