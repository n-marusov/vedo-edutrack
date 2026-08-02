# Domain-Driven Design (DDD) Reference

> Source:
> - https://martinfowler.com/bliki/DomainDrivenDesign.html
> - https://martinfowler.com/bliki/BoundedContext.html
> - https://martinfowler.com/bliki/UbiquitousLanguage.html
> - https://martinfowler.com/bliki/DDD_Aggregate.html
> - https://en.wikipedia.org/wiki/Domain-driven_design
> - https://github.com/ddd-crew/ddd-starter-modelling-process
> - https://github.com/ddd-crew/bounded-context-canvas
> - https://learn.microsoft.com/en-us/azure/architecture/microservices/model/domain-analysis
> - https://learn.microsoft.com/en-us/azure/architecture/microservices/model/tactical-ddd
> Created: 2026-08-02
> Updated: 2026-08-02

## Overview

Domain-Driven Design (DDD) is an approach to software development centered on programming a domain model that has a rich understanding of the processes and rules of a domain. The name comes from Eric Evans's 2003 book *Domain-Driven Design: Tackling Complexity in the Heart of Software*, which describes the approach through a catalog of patterns. DDD is particularly suited to complex domains where a lot of often-messy logic needs to be organized.

DDD is predicated on three goals:

1. Placing the project's primary focus on the core domain and domain logic layer;
2. Basing complex designs on a model of the domain;
3. Initiating creative collaboration between technical and domain experts to iteratively refine a conceptual model that addresses particular domain problems.

DDD rejects the idea of a single unified model for an entire large system — "total unification of the domain model for a large system will not be feasible or cost-effective" (Evans). Instead, it divides a large system into bounded contexts, each of which has its own model. Microsoft recommends DDD only for complex domains where the model provides clear benefits in formulating a common understanding of the domain.

The three pillars of DDD are **Ubiquitous Language**, **Strategic Design**, and **Tactical Design**. Although Evans's background is object-oriented, the core notions are conceptual and apply to any programming approach — this is especially true of strategic design.

## Core Concepts

### Ubiquitous Language

The practice of building up a common, rigorous language between developers and domain experts, based on the Domain Model used in the software. The language must be rigorous because software does not cope well with ambiguity. Using the language in conversations with domain experts is an important part of testing it — and hence the domain model. The language (and the model) should evolve as the team's understanding of the domain grows.

> "By using the model-based language pervasively and not being satisfied until it flows, we approach a model that is complete and comprehensible, made up of simple elements that combine to express complex ideas."
> — Eric Evans

- Domain experts should object to terms or structures that are awkward or inadequate to convey domain understanding.
- Developers should watch for ambiguity or inconsistency that will trip up design.
- Each bounded context can have its own ubiquitous language: the same word (e.g., "account") can have different meanings in different contexts.

### Domain and Domain Model

- **Domain**: the subject area to which the user applies a program.
- **Domain model**: a system of abstractions that describes selected aspects of a domain and can be used to solve problems related to that domain. It distills and organizes domain knowledge and establishes a shared language for developers and domain experts.
- The domain layer is one of the common layers in an object-oriented multilayered architecture.

### Bounded Context

A central pattern in DDD's strategic design section, which deals with large models and teams by dividing them into different bounded contexts and being explicit about their interrelationships.

- A model needs to be unified — internally consistent, with no contradictions. As the modeled domain grows, building a single unified model becomes progressively harder: different groups of people use subtly different vocabularies in different parts of a large organization, and precision runs into polysemes (e.g., "meter", "Customer", "Product" meant subtly different things to different parts of an electricity utility).
- Bounded contexts contain both unrelated concepts (a support ticket only exists in the customer support context) and shared concepts (products, customers) — with completely different models of common concepts and mechanisms to map between them for integration.
- Boundary factors: usually the dominant one is human culture (models act as Ubiquitous Language, so a different model is needed when the language changes); also representation differences (e.g., separation between in-memory and relational database models in a single application).
- Relationships between bounded contexts are depicted using a **context map**.

### Subdomains (Problem Space Classification)

DDD classifies subdomains into three categories to prioritize where to invest the most design effort:

| Category | Definition | Design effort |
|----------|-----------|---------------|
| **Core** | Provides competitive advantage; defines the business | Detailed modeling, substantial team investment |
| **Supporting** | Keeps business operational but doesn't differentiate it | Custom development, but not the source of competitive advantage |
| **Generic** | Problems the industry already solved | Prebuilt or standard solutions instead of custom-built systems |

### Context Mapping Patterns

Context mapping identifies and defines the boundaries of different domains or subdomains within a larger system and visualizes how contexts interact and relate. Patterns (Evans):

| Pattern | Description |
|---------|-------------|
| **Partnership** | Forge a partnership between the teams of two contexts; coordinated planning of development and joint management of integration; the teams succeed or fail together |
| **Shared Kernel** | Designate an explicit boundary around a small subset of the domain model that the teams agree to share; keep the kernel small |
| **Customer/Supplier** | Clear customer/supplier relationship between upstream and downstream teams; the teams negotiate the contract between them |
| **Conformist** | Downstream eliminates the complexity of translation by conforming to the upstream model; used when a custom interface for the downstream subsystem isn't likely to happen |
| **Anti-corruption Layer** | An isolating translation layer that provides the functionality of the upstream system in terms of the downstream's own domain model; protects the model from upstream changes |
| **Open Host Service** | A protocol giving access to a subsystem as a set of services; used when one subsystem must integrate with many others, making custom translations infeasible |
| **Published Language** | A well-documented shared language that expresses the necessary domain information as a common medium of communication (e.g., industry data interchange standards; OpenAPI for REST APIs) |
| **Separate Ways** | Two contexts have no integration; each evolves independently |
| **Big Ball of Mud** | A boundary drawn around the entire mess when no real boundaries exist in an existing system |

### Evans Classification (Tactical Building Blocks)

The book introduced classifying objects into Entities, Value Objects, and Service Objects (the "Evans Classification") and the concept of Aggregates. These filled an important gap that eluded both programming languages and diagrammatic notations.

#### Entity

An object defined not by its attributes but by its identity. Characteristics:

- Has a unique identifier in the system, used to look up or retrieve the entity.
- Two instances sharing the same identity represent the same domain concept even if their attributes differ at a given time (a person's name or address can change, but they remain the same individual); two instances with identical attributes but different identities are distinct.
- The identifier isn't always exposed to users; it can be a GUID or a database primary key.
- **Identity strategy**: natural keys (order number, government-issued ID — convey business meaning, recognized across systems) vs surrogate keys (GUIDs — lack business meaning but avoid coupling to external systems). In a microservices architecture, other services reference entities by identifiers, so identity must remain stable and meaningful across service boundaries; it can span multiple bounded contexts and outlive the application.
- **Encapsulate behavior, not only data**: validation, state transitions, and business rules belong inside the entity. When business logic lives outside in service classes, you create an *anemic domain model* — an antipattern that undermines DDD. Example: a `Delivery` entity should contain the logic deciding whether it can be canceled.

#### Value Object

An object defined only by the values of its attributes; it has no conceptual identity. Two value objects with the same attribute values are interchangeable (people exchanging business cards care about the information on the card, not about distinguishing each unique card).

- **Immutable**: to update a value object, create a new instance to replace the old one. Immutable objects can be shared across threads, cached without defensive copying, and reasoned about more easily in distributed systems.
- May include methods that encapsulate domain logic, but those methods must not produce side effects — they return new value objects.
- **Default modeling choice**: prefer value objects; promote a concept to an entity only when you need to track its identity over time (e.g., `Address` is typically a value object, but becomes an entity if the domain must track a specific address record over time for audit purposes).
- Common examples: colors, dates and times, currency amounts, measurements.

#### Aggregate

A cluster of domain objects that can be treated as a single unit (e.g., an order together with its line items).

- One component object is the **aggregate root**. Any references from outside the aggregate should only go to the root; the root thus ensures the integrity of the aggregate as a whole.
- Aggregates are the basic element of data storage/transfer — you request to load or save whole aggregates. **Transactions should not cross aggregate boundaries.**
- An aggregate defines a consistency boundary around one or more entities and models transactional invariants. An aggregate can consist of a single entity without child entities — what makes it an aggregate is the transactional boundary.
- Aggregate design rules (Microsoft):
  1. **Design small aggregates** — include only the data that must remain consistent within a single transaction. (In drone delivery, `Delivery`, `Package`, `Drone`, and `Account` each form separate aggregates because they have independent life cycles; combining them forces unrelated updates to compete for the same locks.)
  2. **Reference other aggregates by identity only** — the `Delivery` aggregate stores a `DroneId` and a `PackageId`, not direct references.
  3. **Use eventual consistency across aggregates** — when a business process spans multiple aggregates, use domain events rather than a single transaction.
- DDD aggregates are domain concepts (order, clinic visit, playlist), not to be confused with collection classes (lists, maps), which are generic. The term "aggregate" is also used differently in UML.

#### Domain Service and Application Service

A service is a stateless object that implements logic that does not naturally belong to an entity or value object. (This meaning is unrelated to "microservice".)

| Type | Role | Example |
|------|------|---------|
| **Domain service** | Encapsulates business rules that span multiple entities or aggregates | `Scheduler` — scheduling logic involving drone availability, delivery windows, route optimization |
| **Application service** | Orchestrates use cases; coordinates calls to domain services and repositories, manages transactions, handles concerns like authentication or notifications; contains no business logic | An API endpoint that receives a delivery request, calls the `Scheduler`, and returns the result |

#### Domain Events and Integration Events

A domain event is something that happened in the past that domain experts care about — "a delivery was canceled" qualifies; "a record was inserted into a table" does not. Aggregates raise domain events after they change state; events are the primary way to coordinate work across aggregate boundaries.

Classification (Yan Cui, *Serverless Architectures on AWS*):

| Type | Scope | Payload |
|------|-------|---------|
| **Domain events** | Within one bounded context; preserve business logic | Light payloads — only the information needed for processing, since listeners are typically within the same service |
| **Integration events** | Across bounded contexts; ensure data consistency throughout the system | More complex payloads with additional attributes — overcommunication to cover differing listener needs |

Integration events are published asynchronously through a message broker after the transaction that originated them commits. Example: when the shipping bounded context completes a delivery, it publishes a `DeliveryCompleted` integration event that the accounts bounded context consumes to trigger invoices.

Example domain events (drone delivery): `DroneStatus` events describing drone location and status (in-flight, landed); `DeliveryTracking` events when the stage of a delivery changes — `DeliveryCreated`, `DeliveryRescheduled`, `DeliveryHeadedToDropoff`, `DeliveryCompleted`.

#### Repository

An object with methods for retrieving domain objects from a data store (e.g., a database). Creation/retrieval is separated from the domain object itself.

#### Factory

An object with methods for directly creating domain objects. Separates creation from the object itself.

#### Module

A DDD pattern that organizes the model; matters when implementing the model (e.g., inside a microservice).

## Design Methodology

### Strategic vs Tactical Phases

- **Strategic DDD** defines the large-scale system structure: domain analysis, subdomains, bounded contexts, context map. It ensures the architecture remains focused on business capabilities.
- **Tactical DDD** provides design patterns used to create the domain model within a bounded context: entities, aggregates, value objects, domain services, domain events.

### Collaborative Modeling Techniques

- **Event Storming** (Alberto Brandolini): a collaborative, workshop-based modeling technique used as a precursor in DDD to identify and understand domain events. Stakeholders, domain experts, and developers work together to visualize the flow of domain events, their causes, and their effects, using color-coded sticky notes (domain events, aggregates, external systems). It aids in discovering subdomains, bounded contexts, and aggregate boundaries by focusing on "what happens" in the domain. A facilitator experienced with EventStorming helps a team see its benefits beyond a superficial level.
- **Domain Storytelling**, **Example Mapping**, **User Journey Mapping**, **User Story Mapping**, **Business Model Canvas**, **Impact Mapping**, **Product Strategy Canvas**, **Wardley Mapping**, **Core Domain Charts**, **Purpose Alignment Model**, **Business Capability Modelling**, **Domain Message Flow Modelling**, **BPMN**, **Sequence Diagrams**, **Quality Storming**, **Design-Level EventStorming**, **Event Modeling**, **Model Exploration Whirlpool** (Evans), **C4 diagrams**, **UML**.
- **Context Mapper**: a DSL and tools for strategic and tactical DDD.

### Related Architectural Patterns

- **CQRS** (Command Query Responsibility Segregation): separates reading data (queries) from writing data (commands), deriving from CQS (Bertrand Meyer). Commands mutate state and are approximately equivalent to method invocation on aggregate roots; queries read state without mutating it. A command handler invokes a method on the aggregate root; the root performs the logic of the operation and either yields a failure response or mutates its own state that can be written to a data store. The handler pulls in infrastructure concerns (saving the aggregate's state, creating transactions). CQRS does not require DDD, but DDD makes the command/query distinction explicit with aggregate roots.
- **Event Sourcing**: entities track their internal state not by direct serialization/ORM, but by reading and committing events to an event store. Combined with CQRS + DDD: the input is a command and the output is one or many events, saved to an event store and often published on a message broker. Events are often persisted based on the version of the aggregate root, enabling optimistic concurrency in distributed systems.
- **Layered architecture**: the domain layer is one of the common layers in an OO multilayered architecture.
- **Onion / Hexagonal architecture**: used when coding the domain model (aligning code to the domain).
- **Naked Objects**: the UI can simply be a reflection of a good enough domain model; this forces the design of a better domain model.
- **MDE/MDA**: compatible with DDD, but intent differs — MDA is about translating a model into code for different technology platforms, while DDD is about defining better domain models. DSM is DDD applied with domain-specific languages.
- **AOP**: makes it easy to factor technical concerns (security, transactions, logging) out of the domain model.
- **Strangler Fig**: used (with Anti-corruption Layer) when an external/legacy system's schema or API threatens to leak into the application.

## Design Algorithm (Design Process)

### DDD Starter Modelling Process (DDD Crew) — 8 Steps

A step-by-step guide for learning and practically applying each aspect of DDD — from orienting around an organisation's business model to coding a domain model. **Not a linear best practice**: DDD is an evolutionary design process that necessitates continuous iteration on all aspects of knowledge and design. On a real project you jump back and forth between steps.

The 8 steps map to four phases (Eduardo da Silva, "Sociotechnical Architecture"): **Align & Understand**, **Strategic Architecture**, **Strategy & Org Design**, **Tactical Architecture**.

| # | Step | Purpose | Recommended starting point | Additional tools | Who to involve |
|---|------|---------|---------------------------|------------------|----------------|
| 1 | **Understand** | Align focus with the organisation's business model, user needs, and short/medium/long-term goals | Business Model Canvas (business perspective), User Story Mapping (user vantage point) | Impact Mapping, Product Strategy Canvas, Wardley Mapping | Builders, people with domain knowledge, product/business strategy, real end users (not only representatives) |
| 2 | **Discover** | Discover the domain visually and collaboratively; the most crucial aspect of DDD — it cannot be skipped; discovery is continuous | EventStorming | Domain Storytelling, Example Mapping, User Journey Mapping, User Story Mapping | Builders, domain knowledge, product/business strategy, customer-need understanding, real end users |
| 3 | **Decompose** | Decompose the domain into sub-domains — loosely-coupled parts of the domain (reduced cognitive load, team autonomy, loose-coupling/high-cohesion discovery) | Carve the event storm into sub-domains and Context Maps | Business Capability Modelling, Design Heuristics, Independent Service Heuristics | Builders, people with domain knowledge |
| 4 | **Strategize** | Map sub-domains strategically to identify core domains (greatest potential for business differentiation or strategic significance); drives quality/rigour decisions and build vs buy vs outsource | Core Domain Charts | Purpose Alignment Model, Wardley Mapping, "Revisiting the Basics of Domain-Driven Design" | Product/business strategy, builders, domain knowledge |
| 5 | **Connect** | Connect sub-domains into a loosely-coupled architecture that fulfills end-to-end business use-cases; challenge the initial design with concrete use-cases to uncover hidden complexity | Domain Message Flow Modelling | BPMN, Process Modelling EventStorming, Sequence Diagrams | Builders, people with domain knowledge |
| 6 | **Organise** | Organise autonomous teams optimised for fast flow and aligned with context boundaries; account for organisational constraints; teams should self-organise | Context Maps (visualising sociotechnical architecture) | Dynamic Reteaming, Explorers/Villagers/Town Planners, Team Topologies | Builders, domain knowledge, product/business strategy |
| 7 | **Define** | Define the roles and responsibilities of each bounded context; make explicit decisions early, collaboratively and visually, considering technical limitations | Bounded Context Canvas | C4 System Context Diagram, Quality Storming | Builders, domain knowledge, people responsible for the product |
| 8 | **Code** | Code the domain model; aligning code to the domain makes it easier to change the code when the domain changes | Aggregate Design Canvas | C4 Component Diagrams, Design-Level EventStorming, Event Modeling, Hexagonal Architecture, Mob Programming, Model Exploration Whirlpool, Onion Architecture, UML | Builders |

**How to adapt the process** (jump/switch steps when needed):

- Start with collaborative modelling (if the team is more comfortable modeling a familiar domain than talking about business strategy).
- Start by assessing the IT landscape (visualize existing architecture / strategic portfolio first to see major constraints — start with step 5).
- Code before confirming architecture and team boundaries (MVP delivery or a domain so complex that a model in code is needed first).
- Repeat steps 2 (Discover)–6 (Organise) before moving to 7 (Define) to explore different decompositions.
- Organise teams before designing contexts (when organisational constraints dominate).
- Blend Define (7) and Code (8): insights from coding a bounded context can change the high-level design.

### Microsoft's DDD Flow for Microservices (4 Steps)

1. **Analyze the business domain** — understand the application's functional requirements; output is an informal domain description. Map all business functions and their connections (whiteboard sketch or event storming) with domain experts, architects, stakeholders. Identify discrete subdomains: closely related functions, key-vs-supporting functions, the dependency graph; identify external integrations (CRM, payment, billing) but don't focus on technologies yet. External systems may leak their schema/API — establish a boundary (Strangler Fig or Anti-corruption Layer).
2. **Define the bounded contexts** — each bounded context contains a domain model that represents a specific subdomain. Different subsystems may need different models of the same real-world entity (a drone for maintenance/repair vs a drone for scheduling). Group functionality based on whether functions share the same domain model; document interactions in a **context map**.
3. **Apply tactical DDD patterns within a bounded context** — define entities, aggregates, domain services (rules: small aggregates, reference-by-identity, eventual consistency).
4. **Identify the microservices** — design a microservice **no smaller than an aggregate and no larger than a bounded context**.

DDD is an iterative, ongoing process — service boundaries don't remain fixed; as the application evolves you may split a service into several smaller ones. Conway's law: use domain analysis to define boundaries, then intentionally align team ownership to those boundaries; if one team owns multiple unrelated bounded contexts, or one bounded context requires coordination across many teams, revisit boundaries or team structure.

### Bounded Context Canvas (DDD Crew)

A collaborative tool for designing and documenting the design of a single bounded context. Fill the canvas in order: **Name** first, then **Purpose**, then any order (outside-in starting with inbound communication, or inside-out starting with business rules and domain language). Sections:

| Section | Content |
|---------|---------|
| **Name** | The name of the context; agreement on naming frames how the context is designed |
| **Purpose** | A few sentences on the why and what in business language, no technical details; may name key actors |
| **Strategic Classification** | Importance: core / supporting / generic domain. Role in business model: revenue generator / engagement creator / compliance enforcer. Evolution (Wardley Maps): genesis / custom built / product / commodity |
| **Domain Roles** | Behaviour characterisation (analysis vs execution context); see Brandolini's Bounded Context Archetypes, Wirfs-Brock's Object Role Stereotypes |
| **Inbound Communication** | Collaborations initiated by others: messages (command / query / event), collaborators, relationship type (see Context Mapping), swimlanes |
| **Outbound Communication** | Collaborations initiated by this context; same message types and notation |
| **Ubiquitous Language** | Key domain terms within the context and their meanings |
| **Business Decisions** | Key business rules and policies within the context |
| **Assumptions** | Make assumptions explicit — design never happens with full knowledge |
| **Verification Metrics** | Metrics for continuous learning (build-measure-learn): CI/CD environments, JIRA, live systems — to verify whether the chosen boundaries are a good fit |
| **Open Questions** | Questions no one in the room can answer; many questions indicate high uncertainty |

**Interface design tips** (the public interface is the context's contract with the rest of the system): are message names coherent with the description? Is each message type optimal (should a command be an event)? Is the interface too big (too many unique message types)? Does the context expose too much of its internals? Do any messages seem like they should belong elsewhere?

### Example: Drone Delivery (Microsoft)

Strategic analysis: subdomains — **core**: Shipping, Drone management; **supporting**: Invoicing, User accounts (per article: user accounts are called out both as generic elsewhere); **generic**: User accounts, Call center. Bounded contexts model the same entity differently (drone maintenance/repair context vs scheduling context that only needs availability + ETA).

Tactical design of the shipping bounded context:

- **Entities**: `Delivery`, `Package`, `Drone`, `Account`, `Confirmation`, `Notification`, `Tag`.
- **Aggregates**: `Delivery`, `Package`, `Drone`, `Account` — each with its own transactional consistency boundary and independent life cycle. `Confirmation` and `Notification` are child entities of `Delivery`; `Tag` is a child entity of `Package`.
- **Value objects**: `Location`, `ETA`, `PackageWeight`, `PackageSize`.
- **`Delivery` aggregate fields**: ID string; OwnerID: REF; Pickup: Location; Drop-off: Location; Packages: REF; Expedited: BOOLEAN; Confirmation: Confirmation; DroneId: REF. References `Account`, `Package`, `Drone` by identity only.
- **Domain events**: `DroneStatus` (location/status during flight); `DeliveryTracking` — `DeliveryCreated`, `DeliveryRescheduled`, `DeliveryHeadedToDropoff`, `DeliveryCompleted`.
- **Domain services**: `Scheduler` (coordinates all steps to schedule or update a delivery), `Supervisor` (monitors the status of each step to detect failures or timeouts) — a variation of the Scheduler Agent Supervisor pattern.

### Bounded Context → Microservice Mapping

- **1:1** — typically ideal: clear boundaries, reduced coupling, independent deployment and scaling.
- **1:N** — one bounded context divided into multiple microservices to address varying scalability or operational needs.
- **N:1** — multiple bounded contexts consolidated into a single microservice for simplicity or to minimize operational overhead.

The choice balances DDD principles with the system's business goals, technical constraints, and operational requirements.

## Best Practices

1. Build the ubiquitous language together with domain experts and use it consistently in conversations, documentation, and code; challenge awkward terms and watch for ambiguity or inconsistency.
2. Align architecture decisions with the business model, user needs, and goals (Understand step) before choosing technologies.
3. Never skip discovery — it is continuous. Practice event storming and other collaborative techniques frequently; there is always more to learn about the domain.
4. Keep bounded contexts aligned with language/culture boundaries. Do not attempt a single unified enterprise-wide model.
5. Classify subdomains (core / supporting / generic) to focus design effort, decide quality/rigour levels, and make build vs buy vs outsource decisions.
6. Draw a context map and explicitly choose relationship patterns between contexts (ACL, Open Host Service + Published Language, etc.).
7. Design small aggregates; reference other aggregates by identity only; transactions do not cross aggregate boundaries; use eventual consistency with domain events across aggregates.
8. Prefer value objects as the default modeling choice; promote to an entity only when identity over time matters.
9. Encapsulate business rules, validation, and state transitions inside entities and aggregate roots — avoid the anemic domain model.
10. Use domain events to coordinate across aggregate boundaries; publish integration events asynchronously after the originating transaction commits.
11. Model in code and evolve the model during the life of the product (evolutionary design); DDD is a natural component of an Extreme Programming approach.
12. For microservices: a service should be no smaller than an aggregate and no larger than a bounded context; align team ownership to context boundaries (Conway's law).
13. Make assumptions explicit, track verification metrics and open questions (Bounded Context Canvas) to feed the build-measure-learn loop.
14. Protect the model from external/legacy systems: use Anti-corruption Layer or Strangler Fig when an external schema or API threatens to leak in.

## Common Pitfalls

- **Single unified model for a large system**: total unification is not feasible or cost-effective; leads to confusion from polysemic concepts ("Customer", "Product", "meter").
- **Anemic domain model**: business logic lives in service classes instead of entities — undermines the core benefit of DDD, which is expressing business rules inside the domain model.
- **Giant aggregates**: including data that does not need to be consistent in a single transaction forces unrelated updates to compete for the same locks (keep `Delivery`, `Package`, `Drone`, `Account` separate).
- **Transactions across aggregate boundaries**: violates aggregate invariants and blocks distributed designs; use domain events and eventual consistency instead.
- **Overgrown context interface**: too many message types, incoherent names, wrong message type (command vs event), exposing internals.
- **External system leakage**: letting an external system's data schema or API leak into the application compromises the architectural design.
- **Treating the modelling process as a rigid linear sequence**: DDD is evolutionary and iterative; the 8-step process is a starting point to reduce cognitive load, not a standardized best practice.
- **Confusing DDD aggregates with collection classes**: aggregates are domain concepts (order, clinic visit, playlist); collections are generic.
- **Applying DDD where it doesn't pay off**: it requires isolation and encapsulation effort; Microsoft recommends it only for complex domains where the model provides clear benefits.
- **Identity that is not stable across contexts**: other services reference entities by identifier, so identity must remain stable and meaningful across service boundaries.

## Version Notes / Books & Sources

- **Evans, Eric (2003)** — *Domain-Driven Design: Tackling Complexity in the Heart of Software*, Addison-Wesley. The canonical source; opens Part IV with Strategic Design.
- **Evans, Eric (2015)** — *Domain-Driven Design Reference: Definitions and Pattern Summaries* (PDF, freely available at domainlanguage.com).
- **Vernon, Vaughn (2013)** — *Implementing Domain-Driven Design*. A good next step after Evans; focuses on strategic design from the outset (Chapter 2: dividing a domain into bounded contexts; Chapter 3: drawing context maps).
- **Khononov, Vlad** — *Learning Domain-Driven Design: Aligning Software Architecture and Business Strategy*. Practical, modern treatment.
- **Millet, Scott & Tune, Nick (2015)** — *Patterns, Principles, and Practices of Domain-Driven Design* (Wrox).
- **Cui, Yan** — *Serverless Architectures on AWS* (Manning) — domain vs integration event classification.
- **Brandolini, Alberto** — Event Storming; Bounded Context Archetypes.
- **Wirfs-Brock, Rebecca** — Object Role Stereotypes.
- **Foote, Brian & Yoder, Joseph (1999)** — *Big Ball of Mud*.
- **DDD Crew (GitHub)** — `ddd-starter-modelling-process`, `bounded-context-canvas`, `aggregate-design-canvas`, `eventstorming-glossary-cheat-sheet`, `ddd-heuristics`; a DDD Kata by SAP teaches EventStorming + Domain Message Flow + Bounded Context Canvas + Aggregate Canvas.
- **Verraes, Mathias & Wirfs-Brock, Rebecca** — subtleties of delineating bounded contexts (history and human relationships as well as domain concepts).
- **Context Mapper** — a DSL and tools for strategic and tactical DDD.
