<a id="us-execute.gap.cascade-impact-visualization"></a>
# US-execute.gap.cascade-impact-visualization: Визуализация каскадного влияния разрыва

@US-execute.gap.cascade-impact-visualization @UC-execute.gap.diagnose-root-cause @execute @viz @P1
Feature: US-execute.gap.cascade-impact-visualization Визуализация каскадного влияния Gap

  Background:
    Given Найдена корневая причина Gap
    Given Онтология связывает модули по рёбрам hasStrictPrerequisite

  Scenario: Визуализация каскадного влияния с числом модулей и предметов
    Given Диагностирована корневая причина по модулю M
    When Система визуализирует каскадное влияние
    Then Система показывает N модулей, затронутых каскадом из-за M
    And Система показывает распределение затронутых модулей по M предметам
    And Визуализация подсвечивает цепочку от корневой причины до проявления Gap

  Scenario: Каскад не затрагивает другие модули
    Given У корневой причины нет зависимых модулей
    When Система визуализирует влияние
    Then Влияние показывается как нулевое
    And Система возвращает сообщение "No cascade impact: module has no dependents"
