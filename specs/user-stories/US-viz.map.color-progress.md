<a id="us-viz.map.color-progress"></a>
# US-viz.map.color-progress: Цветовое кодирование состояний концептов

@US-viz.map.color-progress @UC-viz.map.view-knowledge-graph-with-progress @P0 @viz
Feature: US-viz.map.color-progress Цветовое кодирование состояний освоения на карте знаний

  Background:
    Given Learner открыл карту знаний
    Given для каждого концепта известно состояние освоения Learner

  Scenario: Зелёный — концепт освоен
    Given Learner освоил концепт "Матрицы"
    When система отображает карту знаний
    Then концепт "Матрицы" окрашен зелёным

  Scenario: Жёлтый — концепт в процессе изучения
    Given Learner начал изучение концепта "Векторы"
    When система отображает карту знаний
    Then концепт "Векторы" окрашен жёлтым

  Scenario: Синий — концепт доступен для изучения
    Given концепт доступен, но изучение не начато
    When система отображает карту знаний
    Then концепт окрашен синим

  Scenario: Серый — концепт недоступен
    Given у концепта не выполнено строгое предусловие hasStrictPrerequisite
    When система отображает карту знаний
    Then концепт окрашен серым

  Scenario: Красный — нарушено строгое предусловие
    Given Learner осваивает концепт, не освоив его строгое предусловие
    When система отображает карту знаний
    Then концепт окрашен красным
    And система показывает предупреждение "Strict prerequisite is not mastered"
