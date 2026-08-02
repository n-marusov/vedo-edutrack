<a id="us-execute.coverage.deficit-list"></a>
# US-execute.coverage.deficit-list: Список дефицитов покрытия с приоритетами

@US-execute.coverage.deficit-list @UC-execute.coverage.deficit-list-with-priority @execute @P1
Feature: US-execute.coverage.deficit-list Список дефицитов ФГОС-покрытия

  Background:
    Given Рассчитано живое покрытие ФГОС
    Given Каждый дефицит связан с модулями и сроками

  Scenario: Формирование ранжированного списка дефицитов
    Given В покрытии ФГОС обнаружены дефициты
    When Система формирует список дефицитов
    Then Дефициты ранжируются по приоритету: влияние на аттестацию, критичность компетенции, ближайший срок
    And Каждый дефицит содержит перечень недостающих модулей
    And Каждый дефицит содержит оценку влияния на ближайший Checkpoint

  Scenario: Дефициты отсутствуют
    Given Покрытие ФГОС полное по всем компетенциям
    When Система формирует список дефицитов
    Then Список дефицитов пуст
    And Система возвращает сообщение "No deficits found"
