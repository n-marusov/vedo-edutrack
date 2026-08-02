<a id="us-execute.alert.deviation"></a>
# US-execute.alert.deviation: Уведомление об отклонении от плана

@US-execute.alert.deviation @UC-execute.alert.deviation-alert @execute @P1
Feature: US-execute.alert.deviation Уведомление при превышении порога отклонения

  Background:
    Given Задан порог отклонения N дней
    Given Событие plan.deviated публикуется при значимом отклонении

  Scenario: Отклонение превышает порог — уведомления Родителю, HR и Learner
    Given Отклонение факта от плана превышает N дней
    When Система обнаруживает превышение порога и публикует событие plan.deviated
    Then Система отправляет уведомление Родителю
    And Система отправляет уведомление HR (при корпоративном обучении)
    And Система отправляет уведомление Learner
    And Уведомление содержит модуль, величину отклонения и рекомендуемое действие

  Scenario: Отклонение в пределах порога
    Given Отклонение факта от плана не превышает N дней
    When Система проверяет отклонение
    Then Уведомления не отправляются
    And Система возвращает сообщение "Deviation within threshold: no alert"
