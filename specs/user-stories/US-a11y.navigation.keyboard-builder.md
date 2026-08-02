<a id="us-a11y.navigation.keyboard-builder"></a>
# US-a11y.navigation.keyboard-builder: Управление конструктором Route с клавиатуры

@US-a11y.navigation.keyboard-builder @UC-a11y.navigation.keyboard-route-builder @P1 @accessibility
Feature: US-a11y.navigation.keyboard-builder Навигация в конструкторе Route без мыши по WCAG 2.1 AA

  Background:
    Given конструктор Route открыт
    Given Learner не использует мышь

  Scenario: Полная навигация с клавиатуры
    Given Learner использует только клавиатуру
    When Learner перемещается по элементам конструктора клавишей Tab
    Then все элементы конструктора доступны с клавиатуры
    And видимый фокус сохраняется на активном элементе

  Scenario: Построение Route без мыши
    Given Learner выбрал цель обучения
    When Learner активировал построение Route клавишей Enter
    Then система строит Route
    And фокус переходит на первый Checkpoint построенной Route

  Scenario: Видимость фокуса
    Given Learner перемещается по конструктору
    When активный элемент меняется
    Then система показывает заметный индикатор фокуса
