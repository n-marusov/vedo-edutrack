<a id="us-viz.panel.group-management"></a>
# US-viz.panel.group-management: Панель управления группой

@US-viz.panel.group-management @UC-viz.panel.group-management-panel @P1 @viz
Feature: US-viz.panel.group-management Панель управления группой: мини-карты, ролевые ограничения, аналитика группы

  Background:
    Given Пользователь аутентифицирован с ролью methodologist или teacher
    Given у Пользователя есть группы Learner

  Scenario: Просмотр группы с мини-картами
    Given Пользователь открыл панель группы
    When система отображает группу
    Then панель показывает мини-карту прогресса для каждого Learner
    And показывает статус Route каждого Learner

  Scenario: Ролевые ограничения видимости
    Given Пользователь имеет роль teacher
    When система отображает панель группы
    Then панель показывает только данные Learner, доступных роли teacher
    And скрывает данные, недоступные роли

  Scenario: Аналитика группы
    Given Пользователь открыл аналитику группы
    When система рассчитывает показатели группы
    Then панель показывает средний прогресс группы
    And показывает долю освоенных Checkpoint
    And показывает типовые Gap группы
