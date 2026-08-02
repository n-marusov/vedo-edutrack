<a id="us-plan.constraint.apply-fgos"></a>
# US-plan.constraint.apply-fgos: Учёт требований ФГОС как фильтра плана

@US-plan.constraint.apply-fgos @UC-plan.constraint.apply-checkpoints-and-fgos @plan @P0
Feature: US-plan.constraint.apply-fgos Требования ФГОС как фильтр плана

  Background:
    Given Существуют требования ФГОС для образовательной программы
    Given Онтология связывает модули с ФГОС-компетенциями

  Scenario: Требования ФГОС применяются как фильтр при построении Route
    Given Родитель указывает программу обучения по ФГОС
    When Система применяет требования ФГОС к построению Route
    Then В Route включаются только модули, покрывающие требуемые ФГОС-компетенции
    And Компетенции без покрытия помечаются как Gap
    And Порядок модулей не нарушает рёбра hasStrictPrerequisite

  Scenario: Программа по ФГОС не указана
    Given Для Learner не указана программа обучения по ФГОС
    When Система применяет требования ФГОС
    Then Фильтр не применяется и Route строится без ФГОС-ограничений
    And Система возвращает сообщение "No FGOS program specified: filter skipped"
