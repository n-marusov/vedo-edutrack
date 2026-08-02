<a id="us-plan.constraint.apply-checkpoints"></a>
# US-plan.constraint.apply-checkpoints: Учёт сроков аттестации как ограничений плана

@US-plan.constraint.apply-checkpoints @UC-plan.constraint.apply-checkpoints-and-fgos @plan @P0
Feature: US-plan.constraint.apply-checkpoints Сроки аттестации как ограничения плана

  Background:
    Given Установлены даты Checkpoint аттестации
    Given Дата целевой аттестации зафиксирована

  Scenario: Сроки аттестации применяются как жёсткие ограничения
    Given Система строит Route с учётом дат аттестации
    When Применяются ограничения по Checkpoint
    Then Каждый Checkpoint становится жёсткой датой в плане
    And Route строится так, чтобы промежуточные Checkpoint были достижимы к своим датам
    And При невозможности уложиться в сроки система помечает модули, требующие переноса

  Scenario: Дата целевой аттестации уже прошла
    Given Дата целевой аттестации ранее текущей даты
    When Система применяет ограничения по Checkpoint
    Then Route не строится
    And Система возвращает сообщение "Attestation date is in the past: constraints unsatisfiable"
