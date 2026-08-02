<a id="us-viz.builder.construct-route"></a>
# US-viz.builder.construct-route: Визуальный конструктор Route

@US-viz.builder.construct-route @UC-viz.builder.construct-route-visually @P0 @viz
Feature: US-viz.builder.construct-route Визуальное построение Route: цель, маршрут, оценка времени, исправление

  Background:
    Given Пользователь открыл конструктор Route
    Given конструктор строит Route как функцию от цели и профиля Learner на основе онтологии VEDO Hub

  Scenario: Построение Route от цели
    Given Пользователь указал цель обучения "Освоить машинное обучение"
    When система строит Route
    Then Route показывает последовательность концептов и Checkpoint от текущего уровня Learner до цели

  Scenario: Оценка времени прохождения
    Given Route построена
    When система рассчитывает время прохождения
    Then система показывает оценку времени для каждого этапа
    And показывает суммарную оценку для всей Route

  Scenario: Исправление построенной Route
    Given Route содержит этап, который не подходит Learner
    When Пользователь заменил этап на альтернативный
    Then система пересчитывает Route с учётом замены
    And система проверяет hasStrictPrerequisite для заменённого этапа

  Scenario: Цель недостижима
    Given цель обучения не связана с графом знаний
    When система строит Route
    Then система показывает ошибку "Goal is not reachable from the knowledge graph"
