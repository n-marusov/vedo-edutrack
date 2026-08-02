<a id="us-api.webhooks.plan-deviated"></a>
# US-api.webhooks.plan-deviated: Вебхук plan.deviated

@US-api.webhooks.plan-deviated @UC-api.webhooks.plan-deviated @P1 @api
Feature: US-api.webhooks.plan-deviated Доставка события plan.deviated внешним системам

  Background:
    Given внешняя система подписана на событие plan.deviated
    Given система доставляет вебхуки с подписью и идемпотентностью

  Scenario: Доставка вебхука при отклонении от плана
    Given план обучения Learner отклонился от графика
    When система публикует событие domain "plan.deviated"
    Then система доставляет вебхук plan.deviated подписчику
    And вебхук содержит id Learner и величину отклонения

  Scenario: Повторная доставка не дублирует обработку
    Given подписчик не подтвердил получение первого вебхука
    When система повторно доставляет вебхук
    Then подписчик получает вебхук с тем же event_id
    And обработка события идемпотентна

  Scenario: Недоступный подписчик
    Given подписчик недоступен
    When система пытается доставить вебхук
    Then система повторяет доставку с экспоненциальной задержкой
    And после исчерпания попыток помечает вебхук как failed
