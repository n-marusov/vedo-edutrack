<a id="us-resource.availability.check-alternatives"></a>
# US-resource.availability.check-alternatives: Проверка доступности ресурса и предложение альтернатив

@US-resource.availability.check-alternatives @UC-resource.availability.check-availability-and-alternatives @P1 @resource
Feature: US-resource.availability.check-alternatives Проверка доступности ресурса и предложение альтернатив

  Background:
    Given ресурс включён в Route
    Given система отслеживает доступность ресурсов по статусу, региону и подписке

  Scenario: Ресурс доступен
    Given Learner открыл ресурс в составе Route
    When Learner запросил проверку доступности
    Then система подтверждает, что ресурс доступен
    And показывает ссылку для перехода к ресурсу

  Scenario: Ресурс недоступен, предложены альтернативы
    Given ресурс недоступен в регионе Learner
    When Learner запросил проверку доступности
    Then система показывает сообщение "Resource is not available in your region"
    And предлагает альтернативные ресурсы, покрывающие ту же цель обучения

  Scenario: Альтернатив нет
    Given ресурс недоступен
    And в онтологии нет эквивалентных ресурсов для цели обучения
    When Learner запросил проверку доступности
    Then система показывает сообщение "No alternatives found"
    And предлагает обратиться к методисту
