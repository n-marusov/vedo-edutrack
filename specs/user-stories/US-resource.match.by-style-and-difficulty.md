<a id="us-resource.match.by-style-and-difficulty"></a>
# US-resource.match.by-style-and-difficulty: Подбор ресурсов по стилю, сложности, длительности и бюджету

@US-resource.match.by-style-and-difficulty @UC-resource.match.match-resources-to-learner @P2 @resource
Feature: US-resource.match.by-style-and-difficulty Подбор ресурсов под профиль Learner (этап 2)

  Background:
    Given профиль Learner содержит стиль обучения, уровень сложности и доступное время
    Given в системе задан бюджет на обучение
    Given каталог содержит ресурсы, размеченные по сложности, длительности и стоимости

  Scenario: Подбор ресурсов по стилю и сложности
    Given Learner предпочитает визуальный стиль и уровень сложности intermediate
    When система подбирает ресурсы для темы "Квантовая механика"
    Then предложенные ресурсы соответствуют стилю и уровню сложности Learner
    And длительность каждого ресурса укладывается в доступное время Learner

  Scenario: Учёт бюджета при подборе
    Given бюджет Learner на обучение составляет 2000 рублей
    When система подбирает платные ресурсы
    Then суммарная стоимость подобранных ресурсов не превышает бюджет
    And для каждого платного ресурса показана бесплатная альтернатива, если она существует

  Scenario: Недостаточно данных профиля
    Given профиль Learner не содержит предпочтений по стилю обучения
    When система подбирает ресурсы
    Then система показывает сообщение "Learner profile is incomplete"
    And предлагает завершить заполнение профиля
