<a id="us-plan.compute.filter-essential"></a>
# US-plan.compute.filter-essential: Фильтрация маршрута по обязательному ядру модулей

@US-plan.compute.filter-essential @UC-plan.compute.filter-by-essential-core @plan @P0
Feature: US-plan.compute.filter-essential Включение обязательного ядра в критический путь

  Background:
    Given Learner обучается по школьной программе с целями, привязанными к ФГОС
    Given VEDO Hub размечает модули как essential core

  Scenario: Модули обязательного ядра включаются в критический путь
    Given Система вычислила Route до цели аттестации
    When Применяется фильтр по essential core
    Then Все модули essential core, необходимые для цели, включаются в Route
    And Модули ядра остаются в критическом пути независимо от оценки длительности
    And Порядок модулей ядра не нарушает рёбра hasStrictPrerequisite

  Scenario: Для цели не размечено обязательное ядро
    Given В онтологии нет модулей essential core, связанных с целью
    When Применяется фильтр по essential core
    Then Route остаётся без изменений
    And Система возвращает сообщение "No essential core modules found for the selected goal"
