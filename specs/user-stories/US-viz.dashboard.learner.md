<a id="us-viz.dashboard.learner"></a>
# US-viz.dashboard.learner: Дашборд Learner

@US-viz.dashboard.learner @UC-viz.dashboard.learner-dashboard @P0 @viz
Feature: US-viz.dashboard.learner Просмотр дашборда Learner: позиция, горизонты, план-факт, покрытие ФГОС

  Background:
    Given Learner аутентифицирован в системе
    Given у Learner есть активная Route

  Scenario: Просмотр позиции в Route
    Given Learner открыл свой дашборд
    When система отображает дашборд
    Then дашборд показывает текущую позицию Learner в Route
    And показывает пройденные и предстоящие Checkpoint

  Scenario: Сравнение плана и факта
    Given у Learner есть план обучения с Checkpoint
    When система отображает дашборд
    Then дашборд показывает план-факт по освоенным Checkpoint
    And показывает отклонение от плана, если оно есть

  Scenario: Покрытие ФГОС
    Given Route построена по требованиям ФГОС
    When система отображает дашборд
    Then дашборд показывает покрытие требований ФГОС
    And показывает непокрытые требования

  Scenario: Автоматическое обновление при событии
    Given система получила событие domain "module.mastered"
    When событие обработано
    Then Mastery Learner обновлён
    And дашборд пересчитан автоматически
