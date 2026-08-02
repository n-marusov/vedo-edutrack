<a id="us-plan.trigger.recompute-progress"></a>
# US-plan.trigger.recompute-progress: Пересчёт маршрута при изменении прогресса

@US-plan.trigger.recompute-progress @UC-plan.compute.recompute-on-progress @plan @P0
Feature: US-plan.trigger.recompute-progress Пересчёт Route при освоении модуля

  Background:
    Given Существует активный Route с планом для Learner
    Given Освоение модуля публикует доменное событие module.mastered

  Scenario: Освоение модуля запускает пересчёт Route
    Given Learner осваивает модуль M и публикуется событие module.mastered
    When Система получает событие module.mastered
    Then Система пересчитывает Route от нового состояния Mastery
    And Система автоматически публикует событие route.recalculated
    And Пересчёт каскадно обновляет план, подбор ресурсов и истории/проекты
    And Route сокращается за счёт исключения освоенного модуля M

  Scenario: Событие module.mastered ссылается на неизвестный модуль
    Given Получено событие module.mastered для модуля, отсутствующего в онтологии
    When Система обрабатывает событие
    Then Пересчёт Route не выполняется
    And Система возвращает сообщение "Unknown module: cannot recompute route"
