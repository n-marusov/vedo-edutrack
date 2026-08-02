<a id="us-resource.catalog.filter-by-format"></a>
# US-resource.catalog.filter-by-format: Фильтрация каталога ресурсов по формату

@US-resource.catalog.filter-by-format @UC-resource.catalog.filter-by-format @P0 @resource
Feature: US-resource.catalog.filter-by-format Фильтрация каталога учебных ресурсов по формату и источнику

  Background:
    Given Learner аутентифицирован в системе
    And каталог ресурсов содержит ресурсы форматов video, article, interactive, exercise и источников hub, external
    And каждый ресурс описан метаданными из онтологии VEDO Hub

  Scenario: Фильтрация каталога по формату
    Given Learner открыл каталог ресурсов по теме "Дифференциальные уравнения"
    When Learner выбрал фильтр формата video
    Then каталог показывает только ресурсы формата video
    And каждый показанный ресурс содержит метку формата

  Scenario: Комбинированная фильтрация по формату и источнику
    Given Learner открыл каталог ресурсов
    When Learner выбрал фильтр формата interactive и источника external
    Then каталог показывает только интерактивные ресурсы из внешних источников

  Scenario: Сброс фильтров
    Given Learner применил фильтры по формату и источнику
    When Learner сбросил фильтры
    Then каталог показывает полный список ресурсов без ограничений

  Scenario: Фильтр без результатов
    Given в каталоге нет ресурсов формата laboratory
    When Learner выбрал фильтр формата laboratory
    Then система показывает сообщение "No resources match the selected filters"
    And предлагает сбросить фильтры
