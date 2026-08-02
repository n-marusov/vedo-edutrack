<a id="us-api.webhooks.module-mastered"></a>
# US-api.webhooks.module-mastered: Вебхук module.mastered

@US-api.webhooks.module-mastered @UC-api.webhooks.module-mastered @P1 @api
Feature: US-api.webhooks.module-mastered Доставка события module.mastered внешним системам

  Background:
    Given внешняя система подписана на событие module.mastered
    Given система доставляет вебхуки с подписью и идемпотентностью

  Scenario: Доставка вебхука при освоении модуля
    Given Learner освоил модуль
    When система публикует событие domain "module.mastered"
    Then система доставляет вебхук module.mastered подписчику
    And вебхук содержит id Learner, id модуля и timestamp

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
