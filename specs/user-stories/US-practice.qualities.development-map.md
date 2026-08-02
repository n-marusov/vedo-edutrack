<a id="us-practice.qualities.development-map"></a>
# US-practice.qualities.development-map: Карта развития личностных качеств

@US-practice.qualities.development-map @UC-practice.qualities.development-map @P2 @practice
Feature: US-practice.qualities.development-map Карта развития качеств для покрытия воспитательных программ (этап 2)

  Background:
    Given у Learner есть Route
    Given в онтологии описаны личностные качества и их связь с учебными активностями

  Scenario: Просмотр карты качеств
    Given Learner открыл карту развития качеств
    When система отображает карту
    Then карта показывает развиваемые качества
    And показывает текущий уровень развития каждого качества

  Scenario: Покрытие воспитательной программы
    Given школа задала воспитательную программу с целевыми качествами
    When система рассчитывает покрытие
    Then карта показывает, какие целевые качества покрываются Route
    And показывает непокрытые качества

  Scenario: Рекомендация активности для качества
    Given качество "Критическое мышление" не покрыто Route
    When система подбирает активности
    Then система предлагает активности, развивающие это качество
