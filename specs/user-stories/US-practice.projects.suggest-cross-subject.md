<a id="us-practice.projects.suggest-cross-subject"></a>
# US-practice.projects.suggest-cross-subject: Предложение межпредметных проектов

@US-practice.projects.suggest-cross-subject @UC-practice.projects.suggest-cross-subject-projects @P1 @practice
Feature: US-practice.projects.suggest-cross-subject Предложение идей межпредметных проектов

  Background:
    Given у Learner есть освоенные концепты из разных предметов
    Given в онтологии заданы связи между концептами разных предметов

  Scenario: Предложение межпредметного проекта
    Given Learner освоил концепты "Электричество" и "Функции"
    When система подбирает проекты
    Then система предлагает проект, объединяющий оба предмета
    And проект опирается на освоенные концепты Learner

  Scenario: Проект с недостающими предусловиями
    Given предложенный проект требует концепт, который Learner не освоил
    When система подбирает проекты
    Then система помечает недостающий концепт
    And показывает связанный Gap

  Scenario: Нет межпредметных связей
    Given освоенные концепты Learner не связаны с проектами
    When система подбирает проекты
    Then система показывает сообщение "No cross-subject projects found"
