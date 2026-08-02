<a id="us-resource.budget.estimate-cost"></a>
# US-resource.budget.estimate-cost: Оценка бюджета Route и сравнение с планом

@US-resource.budget.estimate-cost @UC-resource.budget.estimate-route-budget @P1 @resource
Feature: US-resource.budget.estimate-cost Оценка стоимости Route и сравнение с запланированным бюджетом

  Background:
    Given Learner имеет запланированный бюджет на обучение
    Given Route содержит набор ресурсов с известной стоимостью

  Scenario: Оценка стоимости Route
    Given Route построена для цели "Освоить Python"
    When система рассчитывает бюджет Route
    Then система показывает суммарную стоимость всех ресурсов Route
    And показывает разбивку стоимости по этапам Route

  Scenario: Бюджет в пределах плана
    Given стоимость Route не превышает запланированный бюджет
    When система рассчитывает бюджет Route
    Then система подтверждает, что Route укладывается в запланированный бюджет

  Scenario: Превышение запланированного бюджета
    Given запланированный бюджет Learner составляет 3000 рублей
    And стоимость Route составляет 4500 рублей
    When система рассчитывает бюджет Route
    Then система показывает предупреждение "Route budget exceeds the planned budget"
    And предлагает заменить платные ресурсы на бесплатные альтернативы

  Scenario: Пересчёт после замены ресурса
    Given стоимость Route превышает запланированный бюджет
    And Learner заменил дорогой ресурс на альтернативу
    When система пересчитывает бюджет Route
    Then новая оценка отражает замену
