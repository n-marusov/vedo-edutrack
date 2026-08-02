<a id="us-plan.resource.match-to-step"></a>
# US-plan.resource.match-to-step: Подбор ресурсов к шагам маршрута

@US-plan.resource.match-to-step @UC-plan.resource.match-resources-to-steps @plan @resource @P1
Feature: US-plan.resource.match-to-step Подбор ресурсов к шагам Route

  Background:
    Given Существует Route, построенный до цели аттестации
    Given VEDO Hub предоставляет каталог ресурсов с форматами

  Scenario: Подбор ресурсов подходящего формата к каждому шагу
    Given Route построен до цели аттестации
    When Система подбирает ресурсы к шагам Route
    Then Каждому шагу сопоставляются ресурсы подходящего формата (видео, текст, практика)
    And Ресурсы соответствуют теме шага и уровню Mastery Learner
    And Шаги без подходящих ресурсов помечаются как непокрытые

  Scenario: Для шага нет ресурсов подходящего формата
    Given Для шага отсутствуют ресурсы подходящего формата
    When Система выполняет подбор ресурсов
    Then Шаг помечается как непокрытый ресурсами
    And Система возвращает сообщение "No resources match the required format"
