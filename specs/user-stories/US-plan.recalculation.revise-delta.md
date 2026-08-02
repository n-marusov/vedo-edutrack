<a id="us-plan.recalculation.revise-delta"></a>
# US-plan.recalculation.revise-delta: Пересмотр плана при значимом отклонении

@US-plan.recalculation.revise-delta @UC-plan.recalculation.revise-plan-on-delta @plan @P1
Feature: US-plan.recalculation.revise-delta Предложение пересмотра плана при отклонении

  Background:
    Given Существует зафиксированный план со сроками
    Given Система сравнивает фактический прогресс с планом

  Scenario: Отклонение превышает порог и система предлагает пересмотр
    Given Отклонение факта от плана превышает 15% модулей или 2 недели по срокам
    When Система обнаруживает значимое отклонение и публикует событие plan.deviated
    Then Система формирует предложение пересмотра плана
    And Предложение не применяется автоматически
    And Решение о применении пересмотра принимает пользователь (Родитель или Learner)
    And При согласии пользователя система публикует событие route.recalculated

  Scenario: Отклонение ниже порога
    Given Отклонение факта от плана меньше 15% модулей и меньше 2 недель
    When Система оценивает отклонение
    Then Предложение пересмотра не формируется
    And Система возвращает сообщение "Deviation within tolerance: revision not required"
