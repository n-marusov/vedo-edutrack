<a id="us-execute.gap.diagnose-root-cause"></a>
# US-execute.gap.diagnose-root-cause: Диагностика корневой причины разрыва

@US-execute.gap.diagnose-root-cause @UC-execute.gap.diagnose-root-cause @execute @P0
Feature: US-execute.gap.diagnose-root-cause Диагностика корневой причины Gap

  Background:
    Given Обнаружен Gap в знаниях Learner
    Given Онтология содержит рёбра hasStrictPrerequisite

  Scenario: Поиск корневой причины подъёмом по строгим связям
    Given У Learner обнаружен Gap по модулю M
    When Система поднимается по рёбрам hasStrictPrerequisite от модуля M
    Then Система находит корневой модуль-причину, не освоенный Learner
    And Система показывает цепочку от корневой причины до проявления Gap
    And Система показывает каскадное влияние корневой причины на последующие модули

  Scenario: Цепочка предусловий пуста
    Given Цепочка hasStrictPrerequisite для модуля M пуста
    When Система выполняет диагностику
    Then Корневая причина не определяется
    And Система возвращает сообщение "No root cause found: prerequisite chain is empty"
