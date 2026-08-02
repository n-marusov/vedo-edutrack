<a id="us-execute.coverage.real-vs-formal"></a>
# US-execute.coverage.real-vs-formal: Сопоставление реальных знаний и формальных требований

@US-execute.coverage.real-vs-formal @UC-execute.coverage.real-vs-formal-knowledge-mapping @execute @P1
Feature: US-execute.coverage.real-vs-formal Сопоставление реальных знаний с ФГОС-формулировками

  Background:
    Given Зафиксированы реальные знания Learner по модулям
    Given ФГОС содержит формальные формулировки компетенций

  Scenario: Сопоставление реальных знаний с ФГОС-компетенциями
    Given Известны реальные знания Learner (результаты заданий и практик)
    When Система сопоставляет реальные знания с формальными ФГОС-компетенциями
    Then Для каждой компетенции система показывает, какими реальными знаниями она покрыта
    And Система помечает компетенции, формально покрытые, но фактически неподтверждённые
    And Сопоставление обновляется при появлении новых результатов

  Scenario: Реальные знания не зафиксированы
    Given Нет зафиксированных результатов по Learner
    When Система выполняет сопоставление
    Then Сопоставление не выполняется
    And Система возвращает сообщение "No real knowledge data available for mapping"
