<a id="us-api.sparql.read-only"></a>
# US-api.sparql.read-only: Read-only SPARQL-эндпоинт

@US-api.sparql.read-only @UC-api.sparql.read-only @P1 @api
Feature: US-api.sparql.read-only Read-only SPARQL-доступ к знаниям с аутентификацией и ограничением частоты

  Background:
    Given внешняя система аутентифицирована
    Given SPARQL-эндпоинт работает в режиме read-only

  Scenario: Выполнение read-only SPARQL-запроса
    Given аутентифицированная система отправила SELECT-запрос
    When система выполняет запрос
    Then система возвращает результат запроса
    And результат не позволяет изменять онтологию

  Scenario: Запрос на изменение отклонён
    Given система отправила запрос на изменение данных
    When система обрабатывает запрос
    Then система отклоняет запрос с ошибкой "Method not allowed: endpoint is read-only"

  Scenario: Превышение лимита запросов
    Given система исчерпала лимит запросов за период
    When система отправила следующий запрос
    Then система отклоняет запрос с ошибкой "Rate limit exceeded"
