<a id="us-practice.stories.recommend-at-mastery"></a>
# US-practice.stories.recommend-at-mastery: Рекомендация историй при достижении Mastery

@US-practice.stories.recommend-at-mastery @UC-practice.stories.recommend-stories-at-mastery @P0 @practice
Feature: US-practice.stories.recommend-at-mastery Рекомендация историй при достижении Mastery через appliesTo и enriches

  Background:
    Given в онтологии заданы отношения appliesTo и enriches для историй
    Given у Learner есть освоенные концепты с известным уровнем Mastery

  Scenario: Рекомендация истории после освоения концепта
    Given Learner достиг Mastery по концепту "Алгоритмы"
    When система подбирает истории
    Then система рекомендует истории, связанные с концептом отношением appliesTo
    And рекомендации показывают, на какой концепт они опираются

  Scenario: История, обогащающая концепт
    Given история связана с концептом отношением enriches
    When Learner достиг Mastery по связанному концепту
    Then система предлагает обогащающую историю как дополнительную

  Scenario: Нет подходящих историй
    Given для концепта нет историй с отношениями appliesTo или enriches
    When система подбирает истории
    Then система показывает сообщение "No stories available for this concept"
