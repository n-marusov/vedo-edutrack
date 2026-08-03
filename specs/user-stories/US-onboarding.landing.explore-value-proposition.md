<a id="us-onboarding.landing.explore-value-proposition"></a>
# US-onboarding.landing.explore-value-proposition: Знакомство с продуктом на посадочной странице

@US-onboarding.landing.explore-value-proposition @UC-onboarding.landing.explore-value-proposition @P1 @onboarding
Feature: US-onboarding.landing.explore-value-proposition Посадочная страница как публичная витрина продукта

  Background:
    Given Посетитель не авторизован

  Scenario: Открытие посадочной страницы без авторизации
    When Посетитель открывает корневой URL сервиса
    Then система показывает посадочную страницу без запроса авторизации
    And hero-секция содержит название «VEDO EduTrack» и слоган о персональной траектории

  Scenario: Отображение ценностных предложений
    Then система показывает 3 карточки ценностных предложений: Knowledge Graph Routes, Progress Tracking, Gap Diagnosis

  Scenario: Переход к входу через CTA «Sign In»
    When Посетитель нажимает CTA «Sign In»
    Then система направляет его на страницу входа

  Scenario: Прокрутка к секции возможностей через CTA «Learn More»
    When Посетитель нажимает CTA «Learn More»
    Then система прокручивает страницу к секции возможностей (якорь)

  Scenario: Перенаправление авторизованного пользователя
    Given Посетитель авторизован
    When Посетитель открывает корневой URL сервиса
    Then система перенаправляет его на ролевой дашборд

  Scenario: Локализация RU + EN
    Given Посетитель выбрал английскую локаль
    When Посетитель открывает посадочную страницу
    Then все пользовательские строки отображаются на английском
    And недостающие переводы показываются на русском (fallback)
