<a id="us-plan.compute.cascade-recalculate"></a>
# US-plan.compute.cascade-recalculate: Каскадный пересчёт Route, плана, ресурсов и историй

@US-plan.compute.cascade-recalculate @UC-plan.compute.recompute-on-progress @plan @resource @P1
Feature: US-plan.compute.cascade-recalculate Каскадный пересчёт после пересчёта Route

  Background:
    Given Существует активный Route с планом, ресурсами и историями/проектами
    Given Событие route.recalculated запускает каскадные обновления

  Scenario: Каскад обновляет план, ресурсы и истории/проекты
    Given Система публикует событие route.recalculated
    When Каскад обрабатывает событие
    Then План перестраивается в соответствии с новым Route
    And Подбор ресурсов пересчитывается для новых шагов
    And Истории и проекты (stories/projects) синхронизируются с новым планом
    And Пользователь видит актуальные данные без ручного обновления

  Scenario: Для нового шага каскада нет ресурсов
    Given При каскаде для нового шага отсутствуют ресурсы
    When Каскад обновляет план
    Then Шаг включается в план с пометкой непокрытости ресурсами
    And Система возвращает сообщение "Step has no matched resources"
