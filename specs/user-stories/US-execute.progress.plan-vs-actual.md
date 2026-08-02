<a id="us-execute.progress.plan-vs-actual"></a>
# US-execute.progress.plan-vs-actual: Сравнение плановых и фактических дат

@US-execute.progress.plan-vs-actual @UC-execute.progress.plan-vs-actual-comparison @execute @P0
Feature: US-execute.progress.plan-vs-actual Сравнение план-факт по прогрессу

  Background:
    Given Существует зафиксированный план со сроками
    Given По каждому модулю фиксируются фактические даты освоения

  Scenario: Сравнение плановых и фактических дат с отклонением в днях
    Given Learner завершил модуль M
    When Система сравнивает план-факт по модулю M
    Then Система показывает плановую и фактическую даты модуля M
    And Система вычисляет отклонение в днях: положительное — опережение, отрицательное — отставание
    And Система указывает причину отклонения из зафиксированных фактов

  Scenario: Фактическая дата по модулю не зафиксирована
    Given По модулю M не зафиксирована фактическая дата
    When Система выполняет сравнение план-факт
    Then Отклонение для модуля M не вычисляется
    And Система возвращает сообщение "Actual date not recorded: deviation unknown"
