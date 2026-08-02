<a id="us-execute.assessment.generate-item"></a>
# US-execute.assessment.generate-item: Генерация оценочного задания

@US-execute.assessment.generate-item @UC-execute.assessment.assessment-item-generation @execute @practice @P1
Feature: US-execute.assessment.generate-item Генерация оценочного задания по модулю

  Background:
    Given Модуль связан с компетенциями и критериями Mastery
    Given Существует шаблон генерации оценочных заданий

  Scenario: Генерация задания по модулю и целевому уровню Mastery
    Given Педагог запрашивает оценочное задание по модулю M
    When Система генерирует оценочное задание
    Then Задание соответствует компетенциям модуля M
    And Сложность задания соответствует целевому уровню Mastery
    And Задание сопровождается эталоном оценивания

  Scenario: Для модуля нет шаблона генерации
    Given Для модуля M отсутствует шаблон генерации
    When Система генерирует оценочное задание
    Then Задание не создаётся
    And Система возвращает сообщение "No generation template for module"
