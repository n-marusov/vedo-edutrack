<a id="us-viz.dashboard.parent-hr"></a>
# US-viz.dashboard.parent-hr: Дашборд родителя и HR-специалиста

@US-viz.dashboard.parent-hr @UC-viz.dashboard.parent-hr-dashboard @P1 @viz
Feature: US-viz.dashboard.parent-hr Просмотр дашборда родителя и HR-специалиста

  Background:
    Given Пользователь аутентифицирован с ролью parent или hr
    Given Пользователю назначен Learner для наблюдения

  Scenario: Просмотр прогресса Learner
    Given Пользователь открыл дашборд Learner
    When система отображает дашборд
    Then дашборд показывает прогресс Learner по Route
    And показывает освоенные Checkpoint и уровень Mastery

  Scenario: Просмотр развития компетенций HR-специалистом
    Given HR-специалист наблюдает за Learner
    When система отображает дашборд
    Then дашборд показывает развитие целевых компетенций Learner

  Scenario: Нет доступа к профилю Learner
    Given Пользователь не имеет доступа к выбранному Learner
    When Пользователь запросил дашборд
    Then система показывает ошибку "Forbidden: no access to this learner profile"
