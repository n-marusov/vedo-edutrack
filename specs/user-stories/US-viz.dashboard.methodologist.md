<a id="us-viz.dashboard.methodologist"></a>
# US-viz.dashboard.methodologist: Дашборд методиста и школы

@US-viz.dashboard.methodologist @UC-viz.dashboard.methodologist-dashboard @P1 @viz
Feature: US-viz.dashboard.methodologist Просмотр дашборда методиста и школы

  Background:
    Given Пользователь аутентифицирован с ролью methodologist
    Given методист имеет доступ к группе Learner

  Scenario: Просмотр успеваемости группы
    Given методист открыл дашборд группы
    When система отображает дашборд
    Then дашборд показывает агрегированный прогресс группы
    And показывает распределение уровня Mastery по Learner

  Scenario: Выявление типовых пробелов
    Given методист открыл дашборд группы
    When система отображает дашборд
    Then дашборд показывает типовые Gap, встречающиеся у нескольких Learner
    And показывает предложения по корректировке материалов

  Scenario: Просмотр покрытия ФГОС по школе
    Given методист открыл школьный дашборд
    When система рассчитывает покрытие ФГОС
    Then дашборд показывает покрытие требований по классам и предметам
