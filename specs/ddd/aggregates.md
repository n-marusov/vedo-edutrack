# Агрегаты, сущности и объекты-значения — VEDO EduTrack

> Документ моделирует агрегаты (границы транзакционности), корневые сущности, внутренние сущности и объекты-значения (value objects) для каждого ограниченного контекста из `specs/ddd/context-map.md`. Первичные источники: `specs/glossary.md` (агрегаты `Learner`, `LearningPlan`, `Account`; доменные сервисы `TrajectoryService`, `GapDiagnosisService`, `CoverageService`, `ForecastService`, `PlanService`), `specs/vision.md` §2.2 и §2.4.

## Принципы моделирования

- **Маршрут (Route) — не агрегат**: проекция (производная величина), пересчитываемая по триггерам. Агрегатом является только зафиксированный **план обучения (LearningPlan)** — единственная сущность, которая «берёт на себя обязательство».
- **Траектория (Trajectory) — не агрегат**: производная от истории событий `ModuleMastered`, всегда показывается рядом с маршрутом и планом.
- **Агрегат = граница транзакционности**: внешние ссылки на агрегат только по ID; инварианты агрегата не нарушаются извне.
- **Модули знаний — внешние ссылки**: `KnowledgeModule`, связи, рамки живут в VEDO Hub; в агрегатах EduTrack — только ссылки (ID + версия онтологии).
- **Два контура** (Community / Enterprise) разделяют одну модель; различаются роли, границы видимости и тиры — не структура агрегатов.

---

## 1. Identity & Access (`identity-access`)

### Account (Агрегат)
**Корень:** `Account`
- `accountId: AccountId` (VO)
- `ownerRole: Role` (VO — родитель / директор школы / HR / платформа)
- `tier: Tier` (VO — Community / Pro / Enterprise)
- `status: AccountStatus` (VO)
- `children: List<LearnerRef>` (внутренняя сущность — ссылка на `Learner`)
- `settings: AccountSettings` (VO — границы видимости, тиры, лимиты)

**Инварианты:**
- Родитель видит только своих детей; директор — свою школу; HR — свой департамент (границы видимости F4.7).
- Тировые лимиты (например, до 3 детей на аккаунт в Community) не превышаются.

### Learner (Агрегат)
**Корень:** `Learner`
- `learnerId: LearnerId` (VO)
- `identity: PersonIdentity` (VO — ФИО, контакты)
- `position: Position` (VO) — множество освоений, **источник истины для всех расчётов**
- `goals: List<Goal>` (внутренняя сущность)
- `activePlanId: PlanId | null` (внешняя ссылка на `LearningPlan`)
- `preferences: LearnerPreferences` (VO — темп, стиль обучения, концепция)

**Инварианты:**
- Позиция (`Position`) — единственный источник истины для доступных модулей, маршрута, покрытия, лакун.
- Цель (`Goal`) имеет тип: `role / exam / project / topic / rank / level / permit` (расширяемо).
- Все доменные сервисы принимают позицию как вход (не изменяют её).

### Value Objects (контекст)
`AccountId`, `LearnerId`, `Role`, `Tier`, `AccountStatus`, `PersonIdentity`, `Position`, `Mastery`, `Goal`, `LearnerPreferences`, `AccountSettings`

### Сущности
`Account`, `Learner` (корни), `Goal` (внутренняя сущность `Learner`)

### Доменные события
`GoalChanged` (источник: интерфейс/API — смена цели)

---

## 2. Ontology Port — ACL (`ontology-port`)

> Технический контекст: не владеет предметными агрегатами. Единственный канал чтения онтологии из VEDO Hub.

### Онтологические артефакты (внешние, в VEDO Hub)
`KnowledgeModule`, `Subject`, `LinkType`, `EssentialCore`, `LearningContext`, `PedagogyConcept`, `Story`, `ProjectIdea`, `Quality`, `Standard` (формальная рамка), `OntologyVersion`

### Сущности контекста (порт)
**Корень:** `OntologySnapshot` (кэш-проекция подграфа)
- `snapshotId: SnapshotId` (VO)
- `ontologyVersion: OntologyVersion` (VO) — версия онтологии, на которой вычислен снэпшот
- `subgraph: Graph` (VO — модули, связи, атрибуты)
- `capturedAt: DateTime` (VO)

**Инварианты:**
- Read-only: EduTrack никогда не редактирует онтологию.
- Подграф копируется для вычислений в памяти EduTrack (F0.2).
- Каждый расчёт фиксирует версию онтологии (воспроизводимость).

### Доменные события
Нет (порт — технический; получает `OntologyUpdated`-уведомления от Hub → каскадный пересчёт)

---

## 3. Route Planning (`route-planning`)

### Route (проекция, НЕ агрегат)
- `routeId: RouteId` (VO)
- `position: Position` (VO, вход)
- `goal: Goal` (VO, вход)
- `pedagogyConcept: PedagogyConceptRef` (VO, вход)
- `ontologyVersion: OntologyVersion` (VO, вход)
- `steps: List<RouteStep>` (VO) — шаги маршрута
- `horizons: Horizons` (VO — дальний / средний / ближний)
- `criticalPath: List<RouteStep>` (VO)
- `computedAt: DateTime` (VO)

**Формула:** `Route = f(position, goal, pedagogyConcept, ontologyVersion) → route` — функция, не документ.

### RouteStep (VO)
- `moduleRef: ModuleRef` (VO — внешняя ссылка на `KnowledgeModule`)
- `depth: Depth` (VO) — глубина/контекст прохождения (спиральное обучение: модуль повторяется с возрастающей глубиной)
- `weight: LinkWeight` (VO — strict / soft / enrich)

### Сервисы домена
`RouteComputationService` (двухэтапное вычисление: копирование подграфа → алгоритм маршрутизации), `HorizonService`, `GapToGoalAnalysisService`, `PedagogyConceptService`, `CascadeRecomputeService`

### Доменные события
`RouteRecalculated` (источник: движок EduTrack)

---

## 4. Plan Management (`plan-management`)

### LearningPlan (Агрегат)
**Корень:** `LearningPlan`
- `planId: PlanId` (VO)
- `learnerRef: LearnerId` (внешняя ссылка)
- `routeSnapshot: RouteSnapshot` (VO) — снэпшот маршрута на момент фиксации
- `schedule: Schedule` (внутренняя сущность) — шаги с плановыми датами
- `checkpoints: List<Checkpoint>` (внутренняя сущность)
- `ontologyVersion: OntologyVersion` (VO) — версия онтологии фиксации
- `fixedAt: DateTime` (VO)
- `status: PlanStatus` (VO — активен / пересматривается / завершён)

**Инварианты:**
- План = снэпшот маршрута + временные метки (плановые даты начала/завершения каждого модуля) — фиксируется на контрольной точке.
- После фиксации маршрут продолжает пересчитываться независимо; прогресс измеряется от старого плана, рекомендации — от нового маршрута.
- Пересчёт плана — не автоматический: новый период, смена цели, провал контрольной точки или ручная инициация; дельта >15% модулей или >2 недель → предложение пересмотра (решение пользователя).
- Единственный агрегат, который «берёт на себя обязательство».

### Schedule (внутренняя сущность)
- `items: List<ScheduledStep>` (VO — шаг + плановая дата + нагрузка)
- `constraints: ScheduleConstraints` (VO — нормативные, физиологические, физические, социальные)
- `buffer: Buffer` (VO)

### Checkpoint (внутренняя сущность)
- `date: DateTime` (VO)
- `targetCoverage: CoverageTarget` (VO)
- `status: CheckpointStatus` (VO)

### Сервисы домена
`PlanService` (фиксация, пересмотр), `ScheduleService` (расписание, межпредметная синхронизация F1.10)

### Доменные события
`PlanFixed` (источник: интерфейс)

---

## 5. Execution & Progress (`execution-progress`)

### Trajectory (проекция, НЕ агрегат)
- `learnerRef: LearnerId`
- `masteredEvents: List<MasteryRecord>` (VO) — последовательность освоений с датами
- Производная от событий `ModuleMastered`; не самостоятельный агрегат.

### MasteryRecord (VO)
- `moduleRef: ModuleRef`
- `level: MasteryLevel` (VO — бинарный / градиентный 0–1 / латентный logit)
- `masteredAt: DateTime`
- `source: MasterySource` (VO — тест / подтверждение преподавателя / LMS)

### Deviation (VO)
- `stepRef: RouteStepRef`
- `deviationDays: Integer`
- `reason: DeviationReason` (VO — ускорение / требовалось больше практики / перерыв)

### Сервисы домена
`TrajectoryService` (траектория из событий), `PlanVsActualService` (план-факт), `ForecastService` (прогноз: успевает / под риском / не успевает), `DeviationAlertService` (уведомления)

### Доменные события
`ModuleMastered` (источник: интерфейс/API/LMS), `PlanDeviationDetected` (источник: план обучения)

---

## 6. Gap & Coverage (`gap-coverage`)

### GapDiagnosis (сервис-вычисление, не агрегат)
**Алгоритм:** обнаружено отставание → подъём по графу strict-связей вверх по каждой цепочке пререквизитов → первый неосвоенный модуль на цепочке = корневая лакуна. Несколько цепочек → несколько корневых лакун, ранжируются по каскадному влиянию (число блокируемых модулей и предметов).

### CoverageReport (проекция)
- `coverage: Double` — `coverage = требования рамки, закрытые освоенными модулями / всего требований`
- `deficits: List<Deficit>` (VO) — с приоритетом (жёсткий пререквизит > ядро > опциональное)
- `riskLevel: RiskLevel` (VO — успевает / под риском / не успевает)
- `forecast: Forecast` (VO)

### AssessmentItem (внутренняя сущность, в составе банка)
- `itemId: ItemId` (VO)
- `moduleRef: ModuleRef`
- `difficulty: Difficulty` (VO — IRT/Раша)
- `validity: ValidityIndicator` (VO)
- `reliability: ReliabilityIndicator` (VO)

### Сервисы домена
`GapDiagnosisService`, `CoverageService`, `DeficitService`, `AttestationReadinessService`, `AssessmentService` (генерация + IRT-калибровка)

### Доменные события
`StandardDeficitDetected` (источник: контрольные точки), `AttestationReadinessReportGenerated` (источник: запрос родителя/школы / за 30 дней до КТ)

---

## 7. Resources (`resources`)

### ResourceCatalog (Агрегат)
**Корень:** `ResourceCatalog` (кэш-проекция привязок ресурсов к модулям)
- `entries: List<ResourceEntry>` (внутренняя сущность)

### ResourceEntry (внутренняя сущность)
- `resourceRef: ResourceRef` (VO — ссылка на ресурс, в онтологии/внешнем каталоге)
- `moduleRef: ModuleRef`
- `type: ResourceType` (VO — контент / обеспечение)
- `format: ResourceFormat` (VO — видео, статья, интерактив, учебник, книга)
- `style: LearningStyle` (VO)
- `cost: Cost` (VO)
- `availability: Availability` (VO)

### RouteBudget (проекция)
- `totalCost: Cost` (VO)
- `budgetLimit: Budget` (VO)
- `costDeviations: List<CostDeviation>` (VO)

### Сервисы домена
`ResourceMatchingService` (подбор под ученика), `AvailabilityService` (доступность + альтернативы), `BudgetService` (затраты и бюджет)

### Доменные события
Нет прямых; потребляет `RouteRecalculated` для переподбора ресурсов (F1.6)

---

## 8. Practice & Life (`practice-life`)

### StoryCatalog / ProjectIdeaCatalog (кэш-проекции)
- `stories: List<StoryRef>` (VO — ссылки на истории из онтологии)
- `projectIdeas: List<ProjectIdeaRef>` (VO)
- `qualitiesMap: QualitiesMap` (VO — карта качеств, воспитательная маркировка `develops`)

### Сервисы домена
`StoryRecommendationService` (рекомендация в момент освоения), `ProjectIdeaService` (проектные идеи на стыке модулей), `QualityService` (карта качеств, отбор по целевым качествам)

### Доменные события
`CrossDisciplinaryDiscoveryOffered` (источник: движок EduTrack — межпредметное открытие)

---

## 9. Visualization (`visualization`)

> Read-only контекст: не владеет агрегатами. Потребляет проекции других контекстов (маршрут, горизонты, позиция, лакуны, покрытие, ресурсы).

### View-модели
`KnowledgeMapView`, `GapMapView`, `LearnerDashboardView`, `ParentHrDashboardView`, `MethodologistDashboardView`, `GroupPanelView`, `RouteBuilderView`

### Сервисы домена
`KnowledgeMapProjectionService`, `DashboardProjectionService`, `RouteBuilderService` (конструктор: выбор цели → визуализация → оценка времени → подтверждение → фиксация)

### Доменные события
Нет (только чтение)

---

## 10. Integrations (`integrations`)

> Фасад над ядром. Не владеет агрегатами.

### API-контракты
`RouteApiContract`, `ProgressApiContract`, `CoverageApiContract` — общий для Community и Enterprise.

### Webhook-подписки
`WebhookSubscription` (внутренняя сущность): `endpoint`, `events: List<WebhookEvent>` (VO), `idempotencyKey` (VO), `secret`

### Сервисы домена
`RestApiService`, `SparqlGateway` (read-only, параметризация, rate limiting), `LmsConnectorService` (WebTutor / iSpring / SAP SuccessFactors), `WebhookDispatcherService` (идемпотентная доставка), `McpServerService`, `SsoService` (Keycloak)

### Доменные события
Нет (трансляция событий ядра во внешние webhook-форматы: `route.recalculated`, `module.mastered`, `plan.deviated`, `standard.risk_detected`)

---

## Сводная таблица агрегатов

| Контекст | Агрегаты | Проекции / VO-вычисления | Доменные события |
|----------|----------|--------------------------|------------------|
| Identity & Access | `Account`, `Learner` | `Position`, `Mastery` | `GoalChanged` |
| Ontology Port | `OntologySnapshot` | `Graph`, `OntologyVersion` | — (приём `OntologyUpdated` из Hub) |
| Route Planning | — | `Route`, `RouteStep`, `Horizons`, `CriticalPath` | `RouteRecalculated` |
| Plan Management | `LearningPlan` | `RouteSnapshot`, `Schedule`, `Checkpoint` | `PlanFixed` |
| Execution & Progress | — | `Trajectory`, `MasteryRecord`, `Deviation`, `Forecast` | `ModuleMastered`, `PlanDeviationDetected` |
| Gap & Coverage | — | `GapDiagnosis`, `CoverageReport`, `Deficit`, `AssessmentItem` | `StandardDeficitDetected`, `AttestationReadinessReportGenerated` |
| Resources | `ResourceCatalog` | `ResourceEntry`, `RouteBudget` | — |
| Practice & Life | — | `StoryCatalog`, `ProjectIdeaCatalog`, `QualitiesMap` | `CrossDisciplinaryDiscoveryOffered` |
| Visualization | — | View-модели (read-only) | — |
| Integrations | — | `WebhookSubscription`, API-контракты | — (трансляция во внешние webhook) |

---

## Связанные артефакты

- [Карта контекстов](context-map.md)
- [Каталог доменных событий](domain-events.md)
- [Глоссарий](../glossary.md) — §1, §3, §4
- [Видение продукта](../vision.md) — §2.2, §2.4
