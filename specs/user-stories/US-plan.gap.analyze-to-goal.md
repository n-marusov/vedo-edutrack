<a id="us-plan.gap.analyze-to-goal"></a>
# US-plan.gap.analyze-to-goal: Анализ разрыва до целевой роли

@US-plan.gap.analyze-to-goal @UC-plan.gap.analyze-gap-to-goal @plan @P1
Feature: US-plan.gap.analyze-to-goal Анализ Gap до целевой роли

  Background:
    Given В онтологии существует целевая роль X с набором модулей
    Given Известен текущий Mastery Learner

  Scenario: Ответ на вопрос "Что нужно для роли X?"
    Given Родитель задаёт вопрос "Что нужно для роли X?"
    When Система анализирует Gap между текущим Mastery и требованиями роли X
    Then Система возвращает список недостающих модулей
    And Система выделяет критический путь до роли X
    And Система даёт оценку времени на закрытие Gap
    And Система показывает рёбра hasStrictPrerequisite для недостающих модулей

  Scenario: Роль X отсутствует в онтологии
    Given Роль X не найдена в онтологии
    When Система выполняет анализ Gap
    Then Анализ не выполняется
    And Система возвращает сообщение "Role not found in ontology"
