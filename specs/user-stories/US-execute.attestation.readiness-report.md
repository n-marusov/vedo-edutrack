<a id="us-execute.attestation.readiness-report"></a>
# US-execute.attestation.readiness-report: Отчёт о готовности к аттестации

@US-execute.attestation.readiness-report @UC-execute.attestation.attestation-readiness-report @execute @P1
Feature: US-execute.attestation.readiness-report Отчёт о готовности к аттестации

  Background:
    Given Зафиксирована дата аттестации
    Given Известны Mastery, прогноз готовности и дефициты покрытия

  Scenario: Формирование отчёта о готовности к аттестации
    Given Приближается дата аттестации
    When Система формирует отчёт о готовности
    Then Отчёт содержит бинарный прогноз готовности
    And Отчёт содержит список дефицитов с приоритетами
    And Отчёт содержит отклонения план-факт по ключевым модулям
    And Отчёт содержит рекомендации по закрытию дефицитов до даты аттестации

  Scenario: Дата аттестации не задана
    Given Дата аттестации не зафиксирована
    When Система формирует отчёт
    Then Отчёт не формируется
    And Система возвращает сообщение "Attestation date not set: report unavailable"
