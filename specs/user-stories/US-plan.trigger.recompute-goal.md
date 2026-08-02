<a id="us-plan.trigger.recompute-goal"></a>
# US-plan.trigger.recompute-goal: Полное перестроение маршрута при смене цели

@US-plan.trigger.recompute-goal @UC-plan.compute.recompute-on-goal-change @plan @P1
Feature: US-plan.trigger.recompute-goal Перестроение Route при изменении цели

  Background:
    Given Существует зафиксированный Route для Learner
    Given Цель аттестации может быть изменена Родителем

  Scenario: Смена цели приводит к полному перестроению Route
    Given Родитель изменяет целевую аттестацию для Learner
    When Система получает запрос на пересчёт по смене цели
    Then Система выполняет полное перестроение Route от текущего Mastery к новой цели
    And Старый Route полностью заменяется новым
    And Система публикует событие route.recalculated

  Scenario: Новая цель идентична текущей
    Given Родитель сохраняет цель без изменений
    When Система проверяет запрос на пересчёт
    Then Полное перестроение не выполняется
    And Система возвращает сообщение "Goal unchanged: no rebuild required"
