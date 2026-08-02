<a id="us-api.rest.compute-route"></a>
# US-api.rest.compute-route: REST-эндпоинт вычисления Route

@US-api.rest.compute-route @UC-api.rest.compute-route @P0 @api
Feature: US-api.rest.compute-route Вычисление Route через REST API для EdTech-систем

  Background:
    Given внешняя EdTech-система аутентифицирована по API-ключу
    Given API предоставляет эндпоинт POST /api/v1/routes/compute

  Scenario: Успешное вычисление Route
    Given EdTech-система отправила запрос с целью обучения и профилем Learner
    When система вычисляет Route
    Then API возвращает статус 200
    And ответ содержит Route с последовательностью концептов и Checkpoint

  Scenario: Некорректный запрос
    Given запрос не содержит цели обучения
    When EdTech-система вызывает эндпоинт
    Then API возвращает ошибку "Bad Request: goal is required"

  Scenario: Недостаточные права
    Given EdTech-система не имеет прав на вычисление Route
    When EdTech-система вызывает эндпоинт
    Then API возвращает ошибку "Forbidden"
