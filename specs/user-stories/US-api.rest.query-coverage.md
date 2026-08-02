<a id="us-api.rest.query-coverage"></a>
# US-api.rest.query-coverage: REST-эндпоинт запроса покрытия ФГОС

@US-api.rest.query-coverage @UC-api.rest.query-coverage @P1 @api
Feature: US-api.rest.query-coverage Запрос покрытия требований через REST API

  Background:
    Given внешняя система аутентифицирована по API-ключу
    Given API предоставляет эндпоинт GET /api/v1/learners/{id}/coverage

  Scenario: Запрос покрытия ФГОС
    Given у Learner есть Route, построенная по требованиям ФГОС
    When внешняя система запросила покрытие
    Then API возвращает статус 200
    And ответ содержит покрытые и непокрытые требования

  Scenario: Покрытие с фильтром по предмету
    Given внешняя система запросила покрытие с фильтром по предмету "Математика"
    When система обрабатывает запрос
    Then ответ содержит только требования по указанному предмету

  Scenario: Покрытие недоступно
    Given у Learner нет Route
    When внешняя система запросила покрытие
    Then API возвращает ошибку "No route found for the learner"
