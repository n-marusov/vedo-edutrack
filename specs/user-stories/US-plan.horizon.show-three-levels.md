<a id="us-plan.horizon.show-three-levels"></a>
# US-plan.horizon.show-three-levels: Отображение трёх горизонтов планирования

@US-plan.horizon.show-three-levels @UC-plan.horizon.show-three-horizons @plan @viz @P0
Feature: US-plan.horizon.show-three-levels Отображение трёх горизонтов плана

  Background:
    Given Существует Route до цели аттестации
    Given План разбит на горизонты по длительности

  Scenario: Отображение ближнего, среднего и дальнего горизонтов
    Given Родитель открывает план Learner
    When Система отображает план
    Then Система показывает ближний горизонт (near): ближайшие шаги и ближайший Checkpoint
    And Система показывает средний горизонт (mid): этапы до промежуточной аттестации
    And Система показывает дальний горизонт (far): полный путь до целевой аттестации
    And Горизонты согласованы: дальний включает средний, средний включает ближний

  Scenario: Route короче дальнего горизонта
    Given Длительность Route меньше минимальной границы дальнего горизонта
    When Система отображает план
    Then Дальний горизонт отображается как полный Route без разбиения
    And Система возвращает сообщение "Route shorter than far horizon: no split applied"
