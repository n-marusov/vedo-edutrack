<a id="us-a11y.navigation.screen-reader-map"></a>
# US-a11y.navigation.screen-reader-map: Доступ к карте знаний через скринридер

@US-a11y.navigation.screen-reader-map @UC-a11y.screen-reader.knowledge-map @P0 @accessibility
Feature: US-a11y.navigation.screen-reader-map Доступ к карте знаний через скринридер по WCAG 2.1 AA

  Background:
    Given Learner использует скринридер
    Given карта знаний содержит визуальный граф концептов

  Scenario: Озвучивание концептов карты
    Given Learner открыл карту знаний
    When скринридер читает карту
    Then каждый концепт озвучивается с названием и состоянием освоения

  Scenario: Озвучивание связей между концептами
    Given на карте есть связь с hasStrictPrerequisite
    When скринридер читает связь
    Then скринридер объявляет связь и её тип

  Scenario: Альтернативный список концептов
    Given Learner с ограниченным зрением открыл карту знаний
    When Learner запросил список концептов
    Then система показывает линейный список всех концептов с их состоянием
    And список доступен с клавиатуры
