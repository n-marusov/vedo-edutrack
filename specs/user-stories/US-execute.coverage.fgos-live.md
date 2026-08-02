<a id="us-execute.coverage.fgos-live"></a>
# US-execute.coverage.fgos-live: Живое покрытие требований ФГОС

@US-execute.coverage.fgos-live @UC-execute.coverage.fgos-coverage-live @execute @viz @P0
Feature: US-execute.coverage.fgos-live Живое отображение покрытия ФГОС

  Background:
    Given Программа обучения Learner привязана к ФГОС
    Given Mastery по модулям обновляется событиями module.mastered

  Scenario: Отображение покрытия ФГОС в реальном времени
    Given Родитель открывает дашборд покрытия ФГОС
    When Система отображает покрытие в реальном времени
    Then Система показывает процент покрытия каждой ФГОС-компетенции
    And Покрытие обновляется автоматически при каждом событии module.mastered
    And Компетенции с полным покрытием выделяются отдельно

  Scenario: ФГОС-компетенции не привязаны к модулям
    Given ФГОС-компетенции не связаны с модулями в онтологии
    When Система отображает покрытие
    Then Покрытие не рассчитывается
    And Система возвращает сообщение "FGOS mapping unavailable: coverage cannot be computed"
