<a id="us-plan.fixation.snapshot"></a>
# US-plan.fixation.snapshot: Фиксация плана на контрольной точке

@US-plan.fixation.snapshot @UC-plan.fixation.snapshot-plan @plan @P0
Feature: US-plan.fixation.snapshot Фиксация плана на Checkpoint

  Background:
    Given Существует активный Route для Learner
    Given Проводятся Checkpoint аттестации с фиксированными датами

  Scenario: Фиксация текущего плана на Checkpoint
    Given Наступает дата Checkpoint аттестации
    When Система фиксирует текущий план как снимок (snapshot)
    Then Снимок сохраняется для последующего сравнения план-факт
    And Route продолжает пересчитываться по событиям module.mastered
    And Система отображает два слоя: зафиксированный план и актуальный пересчитанный Route

  Scenario: Фиксация при отсутствии активного плана
    Given Для Learner нет активного Route
    When Система пытается зафиксировать план
    Then Снимок не создаётся
    And Система возвращает сообщение "No active plan to snapshot"
