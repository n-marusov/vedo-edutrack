<a id="us-plan.compute.shortest-path"></a>
# US-plan.compute.shortest-path: Построение кратчайшего маршрута до цели аттестации

@US-plan.compute.shortest-path @UC-plan.compute.shortest-path-to-goal @plan @P0
Feature: US-plan.compute.shortest-path Построение кратчайшего маршрута до цели

  Background:
    Given Learner имеет зафиксированную цель аттестации (Checkpoint)
    Given VEDO Hub предоставляет онтологию с рёбрами hasStrictPrerequisite

  Scenario: Построение кратчайшего маршрута с учётом строгих предусловий
    Given Родитель выбирает цель аттестации для Learner
    When Система вычисляет кратчайший путь от текущего Mastery Learner до цели
    Then Система возвращает Route с упорядоченным набором модулей
    And Route учитывает рёбра hasStrictPrerequisite: модуль не планируется до выполнения строгих предусловий
    And Route не включает модули, по которым Mastery уже достигнут

  Scenario: Цель недостижима из-за строгих предусловий
    Given В онтологии отсутствует выполнимая цепочка hasStrictPrerequisite до цели
    When Система пытается вычислить кратчайший путь
    Then Route не создаётся
    And Система возвращает сообщение "No feasible route: strict prerequisites are unsatisfied"
