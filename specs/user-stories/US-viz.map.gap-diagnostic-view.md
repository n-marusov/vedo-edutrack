<a id="us-viz.map.gap-diagnostic-view"></a>
# US-viz.map.gap-diagnostic-view: Диагностическая карта пробелов (Gap)

@US-viz.map.gap-diagnostic-view @UC-viz.map.view-gap-diagnostic-map @P0 @viz
Feature: US-viz.map.gap-diagnostic-view Просмотр карты пробелов с корневыми Gap и каскадными блокировками

  Background:
    Given у Learner проведена диагностика знаний
    Given диагностика выявила пробелы Gap в знаниях
    Given для Gap определены связи через hasStrictPrerequisite

  Scenario: Просмотр корневых пробелов
    Given Learner открыл диагностическую карту
    When система отображает карту Gap
    Then корневые Gap показаны в начале цепочек
    And каждый Gap показан с текущим уровнем Mastery

  Scenario: Каскадная блокировка последующих концептов
    Given Gap по концепту "Пределы" блокирует концепт "Интегралы" через hasStrictPrerequisite
    When система отображает карту Gap
    Then между Gap и заблокированным концептом показана стрелка блокировки
    And заблокированные концепты отмечены красным

  Scenario: Обновление карты при событии освоения
    Given система получила событие domain "module.mastered" для концепта с Gap
    When событие обработано
    Then Gap закрывается
    And каскадно пересчитываются зависимые концепты

  Scenario: Нет пробелов
    Given диагностика не выявила пробелов
    When Learner открыл диагностическую карту
    Then система показывает сообщение "No knowledge gaps found"
