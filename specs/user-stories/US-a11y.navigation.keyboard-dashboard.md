<a id="us-a11y.navigation.keyboard-dashboard"></a>
# US-a11y.navigation.keyboard-dashboard: Доступ к дашборду с клавиатуры и скринридером

@US-a11y.navigation.keyboard-dashboard @UC-a11y.screen-reader.knowledge-map @P1 @accessibility
Feature: US-a11y.navigation.keyboard-dashboard Доступ к дашборду с клавиатуры и скринридером по WCAG 2.1 AA

  Background:
    Given дашборд Learner открыт
    Given Learner использует клавиатуру или скринридер

  Scenario: Навигация по дашборду с клавиатуры
    Given Learner использует только клавиатуру
    When Learner перемещается по виджетам дашборда
    Then все виджеты дашборда доступны с клавиатуры
    And видимый фокус сохраняется на активном виджете

  Scenario: Чтение данных дашборда скринридером
    Given Learner использует скринридер
    When скринридер читает дашборд
    Then показатели прогресса озвучиваются с заголовками
    And озвучиваются отклонения от плана

  Scenario: Скрытая информация недоступна скринридеру
    Given виджет скрыт с экрана
    When скринридер читает дашборд
    Then система не озвучивает скрытый виджет
