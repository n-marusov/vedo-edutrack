<a id="us-viz.map.knowledge-graph"></a>
# US-viz.map.knowledge-graph: Карта знаний с прогрессом Learner

@US-viz.map.knowledge-graph @UC-viz.map.view-knowledge-graph-with-progress @P0 @viz
Feature: US-viz.map.knowledge-graph Просмотр графа знаний с прогрессом и строгими предусловиями

  Background:
    Given Learner имеет Route в системе
    Given карта знаний построена на основе онтологии VEDO Hub
    Given для концептов заданы отношения hasStrictPrerequisite

  Scenario: Просмотр карты знаний с прогрессом
    Given Learner открыл карту знаний
    When система отображает граф знаний
    Then концепты окрашены по состоянию освоения Learner
    And стрелки показывают связи между концептами

  Scenario: Подсветка нарушенных строгих предусловий
    Given концепт "Интегралы" имеет строгое предусловие hasStrictPrerequisite на концепт "Пределы"
    Given Learner не освоил концепт "Пределы"
    When система отображает карту знаний
    Then концепт "Интегралы" выделен красным
    And система показывает предупреждение "Strict prerequisite is not mastered"

  Scenario: Пустая карта знаний
    Given у Learner нет освоенных концептов
    When система отображает карту знаний
    Then карта показывает концепты без прогресса
    And все концепты имеют нейтральный статус
