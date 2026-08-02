<a id="us-execute.forecast.binary-readiness"></a>
# US-execute.forecast.binary-readiness: Бинарный прогноз готовности к аттестации

@US-execute.forecast.binary-readiness @UC-execute.forecast.binary-readiness-forecast @execute @P0
Feature: US-execute.forecast.binary-readiness Бинарный прогноз готовности к Checkpoint

  Background:
    Given Существует ближайший Checkpoint аттестации
    Given Известны текущий темп освоения и отклонения от плана

  Scenario: Прогноз готовности как один из двух статусов
    Given Наступает момент прогноза до ближайшего Checkpoint
    When Система строит прогноз готовности
    Then Система выдаёт один из статусов: on-track или not-on-track
    And Прогноз учитывает оставшиеся модули, темп и рёбра hasStrictPrerequisite
    And При статусе not-on-track система показывает список рисковых модулей

  Scenario: Недостаточно данных для прогноза
    Given Нет ни одного зафиксированного факта освоения модулей
    When Система строит прогноз
    Then Прогноз не выдаётся
    And Система возвращает сообщение "Insufficient data for readiness forecast"
