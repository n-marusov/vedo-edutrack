<a id="us-api.sso.keycloak"></a>
# US-api.sso.keycloak: Вход через Keycloak SSO

@US-api.sso.keycloak @UC-api.sso.keycloak-sso-integration @P1 @api
Feature: US-api.sso.keycloak Вход в систему через Keycloak SSO для Enterprise

  Background:
    Given организация использует Keycloak как identity provider
    Given система интегрирована с Keycloak через SSO

  Scenario: Успешный вход через Keycloak
    Given Learner открыл страницу входа
    When Learner выбрал вход через корпоративный SSO
    Then система перенаправляет Learner в Keycloak
    And после успешной аутентификации Learner попадает в систему

  Scenario: Новый пользователь из Keycloak
    Given Learner впервые входит через Keycloak
    When Keycloak вернул данные пользователя
    Then система создаёт профиль Learner
    And связывает профиль с учётной записью Keycloak

  Scenario: Отмена входа в Keycloak
    Given Learner отменил вход в Keycloak
    When Learner вернулся в систему
    Then система показывает сообщение "Authentication cancelled"
    And Learner остаётся неаутентифицированным
