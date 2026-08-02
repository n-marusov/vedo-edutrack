# Clean Architecture Reference

> Source: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
> https://learn.microsoft.com/en-us/dotnet/architecture/modern-web-apps-azure/common-web-application-architectures
> https://alistair.cockburn.us/hexagonal-architecture/
> https://jeffreypalermo.com/blog/the-onion-architecture-part-1/
> https://github.com/ardalis/cleanarchitecture
> Created: 2026-08-02
> Updated: 2026-08-02

## Overview

Clean Architecture is a software architecture methodology introduced by Robert C. Martin (Uncle Bob) in his 2012 blog post *The Clean Architecture* and later formalized in the book *Clean Architecture: A Craftsman's Guide to Software Structure and Design* (2017). It is "just the latest in a series of names for the same loosely-coupled, dependency-inverted architecture" — the same idea is known as Hexagonal Architecture / Ports-and-Adapters (Alistair Cockburn, 2005), Onion Architecture (Jeffrey Palermo, 2008), and Screaming Architecture (Martin, 2011), and it descends from DCI (Coplien/Reenskaug) and BCE (Ivar Jacobson).

The core idea is **separation of concerns** achieved by dividing the software into concentric layers: business rules at the center, mechanisms (frameworks, database, UI) on the outside. Systems built this way are:

- **Independent of Frameworks** — frameworks are used as tools, not as constraints the system must be crammed into.
- **Testable** — business rules can be tested without the UI, Database, Web Server, or any other external element.
- **Independent of UI** — a web UI can be replaced with a console UI without changing business rules.
- **Independent of Database** — Oracle/SQL Server can be swapped for Mongo/BigTable/CouchDB without touching business rules.
- **Independent of any external agency** — business rules know nothing about the outside world.

Note: the English Wikipedia article "Clean_architecture" does not exist (404 as of 2026-08-02); the methodology is documented via the primary sources above.

## Core Concepts

### The Dependency Rule

The overriding rule that makes this architecture work. **Source code dependencies can only point inwards.** Nothing in an inner circle can know anything at all about something in an outer circle — the name of anything declared in an outer circle (functions, classes, variables, any named entity) must not be mentioned by inner code. Data formats used in an outer circle (especially framework-generated formats) must not be used by an inner circle.

- Outer circles are **mechanisms**; inner circles are **policies**.
- Moving inward, the level of abstraction increases; the outermost circle is low-level concrete detail, the innermost circle is the most general.

### Entities

Encapsulate **enterprise-wide business rules**. An entity can be an object with methods or a set of data structures and functions — what matters is that entities could be used by many different applications in the enterprise. In a single application, entities are the application's business objects. They are the least likely to change when something external changes: no change to page navigation, security, or any operational change to a particular application should affect the entity layer.

### Use Cases

Contains **application-specific business rules**. Encapsulates and implements all use cases of the system. Use cases orchestrate the flow of data to and from the entities and direct entities to use their enterprise-wide business rules to achieve the goal of the use case. Expected properties:

- Changes in this layer do not affect entities.
- This layer is not affected by changes to externalities (database, UI, common frameworks).
- Changes to the operation of the application **will** affect the use cases layer.

### Interface Adapters

A set of adapters that convert data between the format most convenient for use cases/entities and the format most convenient for an external agency (Database, Web).

- The MVC architecture of a GUI lives wholly in this layer: Presenters, Views, Controllers; models are data structures passed from controllers to use cases and back to presenters/views.
- Data is converted from the form most convenient for entities/use cases to the form required by the persistence framework. If the database is SQL, **all SQL is restricted to this layer** (specifically the DB-related parts).
- Also contains any adapter to convert data from an external service to the internal form used by use cases/entities.
- No code inward of this circle knows anything at all about the database.

### Frameworks and Drivers

The outermost layer: the Database, the Web Framework, UI toolkit, etc. Generally only glue code that communicates to the next circle inward. "This layer is where all the details go. The Web is a detail. The database is a detail. We keep these things on the outside where they can do little harm."

### Only Four Circles?

The circles are schematic. More than four layers may be needed — there is no rule requiring exactly four. **The Dependency Rule always applies** regardless of the number of layers.

### Crossing Boundaries

Flow of control starts in the controller, moves through the use case, and winds up executing in the presenter. The apparent contradiction with the Dependency Rule (control flows outward, dependencies must point inward) is resolved with the **Dependency Inversion Principle**: arrange interfaces and inheritance so that source code dependencies oppose the flow of control at the boundary.

Example: a use case needs to call the presenter, but a direct call would violate the Dependency Rule (no name of an outer circle may be mentioned by an inner circle). Instead, the use case calls an interface declared in the inner circle — the **Use Case Output Port** — and the presenter (outer circle) implements it. Dynamic polymorphism makes source dependencies oppose the flow of control at every boundary.

### Data Crossing Boundaries

Data crossing boundaries is typically simple data structures: basic structs, simple Data Transfer Objects (DTOs), function-call arguments, hashmaps, or objects. Requirements:

- Pass **isolated, simple data structures** — never Entities, never Database rows.
- The data structures must not carry any dependency that violates the Dependency Rule.
- Data is always passed in the form **most convenient for the inner circle**.

Example from Martin: database frameworks return a convenient format ("RowStructure") for query results — that row structure must never be passed inward across a boundary, because it forces an inner circle to know about an outer circle.

## Related Architectures

| Architecture | Author | Year | Core metaphor |
|---|---|---|---|
| Hexagonal (Ports & Adapters) | Alistair Cockburn | 2005 | Hexagon; application inside communicates over **ports** to external agencies; technology-specific **adapters** convert port APIs to device signals |
| Onion | Jeffrey Palermo | 2008 | Concentric layers; Domain Model in the very center; **all coupling toward the center** |
| Screaming | Robert C. Martin | 2011 | Architecture should "scream" about the use cases, not about the frameworks |
| DCI (Data-Context-Interaction) | James Coplien, Trygve Reenskaug | — | Data, Context, and Interaction separation |
| BCE (Boundary-Control-Entity) | Ivar Jacobson | — | Use-case-driven decomposition (Object Oriented Software Engineering) |
| Clean | Robert C. Martin | 2012 | Concentric circles; the Dependency Rule |

### Hexagonal Architecture (Ports and Adapters) — key points

Cockburn's intent: "Create your application to work without either a UI or a database so you can run automated regression-tests against the application, work when the database becomes unavailable, and link applications together without any user involvement."

- A **port** identifies a purposeful conversation (an API); the protocol of a port is given by the purpose of the conversation.
- An **adapter** converts the API definition to the signals needed by a device and vice versa. Examples: GUI, test harness (FIT/Fitnesse), batch driver, HTTP interface, direct program-to-program interface, mock (in-memory) database, real database.
- The rule: **code pertaining to the inside part should not leak into the outside part**.
- Primary (driving) ports/adapters = actors that drive the application (natural test adapter: FIT). Secondary (driven) ports/adapters = what the application drives (natural test adapter: mock).
- Use cases should be written at the application boundary (the inner hexagon), specifying functions/events supported by the application regardless of external technology.
- Number of ports is a matter of taste; 2–4 is typical (e.g., weather system: feed, administrator, subscribers, subscriber DB).

### Onion Architecture — key points

- "All code can depend on layers more central, but code cannot depend on layers further out from the core. In other words, all coupling is toward the center."
- Domain Model at the very center (state + behavior that models truth for the organization); the Domain Model is only coupled to itself.
- The first ring around the domain holds **repository interfaces** (object saving/retrieving) — only the interface is in the application core; the implementation (coupled to a particular data-access method) lives on the edge.
- Outer layer (UI, Infrastructure, Tests) is reserved for things that change often, intentionally isolated from the application core.
- **The database is not the center. It is external.** "With Onion Architecture, there are no database applications."
- Relies heavily on the Dependency Inversion principle; implementations at the edges are injected at runtime.
- Appropriate for long-lived business applications and applications with complex behavior — **not** for small websites.

## Usage Patterns

### Pattern 1: Crossing a boundary (use case → presenter) with the Dependency Rule

Source: Martin's description of the boundary-crossing example (lower-right of his diagram). The generic mechanism:

1. The use case defines an interface in its own (inner) circle — the **Use Case Output Port** — declaring the callback it needs from the outside.
2. The presenter (outer circle) **implements** that interface.
3. At runtime, dependency injection wires the presenter to the use case; the source dependency still points inward (use case → interface), while control flows outward.

The same technique (interface in the inner circle, implementation in the outer circle) is used to cross every boundary in the architecture, exploiting dynamic polymorphism.

### Pattern 2: Ports & Adapters — Cockburn's Discounter example (verbatim, Java)

The simplest ports-and-adapters application: `discount(amount) = amount * rate(amount)`. Amount comes from the user, rate from a database → two ports. Implemented in stages: (1) tests with a constant rate, (2) GUI, (3) mock database swappable for a real database.

The database-side adapter (replaceable repository interface):

```java
public interface RateRepository {
    double getRate(double amount);
}
```

```java
public class RepositoryFactory {
    public RepositoryFactory() { super(); }

    public static RateRepository getMockRateRepository() {
        return new MockRateRepository();
    }
}
```

```java
public class MockRateRepository implements RateRepository {
    public double getRate(double amount) {
        if(amount <= 100) return 0.01;
        if(amount <= 1000) return 0.02;
        return 0.05;
    }
}
```

The application accepts the repository adapter via constructor injection:

```java
import repository.RepositoryFactory;
import repository.RateRepository;

public class Discounter {
    private RateRepository rateRepository;

    public Discounter(RateRepository r) {
        super();
        rateRepository = r;
    }

    public double discount(double amount) {
        double rate = rateRepository.getRate( amount );
        return amount * rate;
    }
}
```

The user-side test adapter (FIT ColumnFixture) passes a mock repository:

```java
import app.Discounter;
import fit.ColumnFixture;

public class TestDiscounter extends ColumnFixture {
    private Discounter app =
        new Discounter(RepositoryFactory.getMockRateRepository());
    public double amount;

    public double discount() {
        return app.discount( amount );
    }
}
```

Development order this enables: FIT test harness + mock DB → add GUI (still on mock DB) → integration tests against a real DB with test data → real use with a person + live database.

### Pattern 3: Microsoft's project organization (.NET / eShopOnWeb / ardalis template)

Clean Architecture puts business logic and the application model at the center; "infrastructure and implementation details depend on the Application Core." Solid arrows = compile-time dependencies; dashed = runtime-only.

**Application Core** (no dependencies on other layers):
- Entities (business model classes that are persisted)
- Aggregates (groups of entities)
- Interfaces (abstractions for operations performed using infrastructure: data access, file system, network calls)
- Domain Services
- Specifications
- Custom Exceptions and Guard Clauses
- Domain Events and Handlers
- DTOs (simple non-entity types with no UI/Infrastructure dependencies)

**Infrastructure** (implements interfaces defined in Application Core; references Application Core):
- EF Core types (`DbContext`, `Migration`)
- Data access implementation types (Repositories)
- Infrastructure-specific services (e.g., `FileLogger`, `SmtpNotifier`)

**UI Layer** (entry point; references Application Core; interacts with Infrastructure strictly through interfaces defined in Application Core — no direct instantiation or static calls to Infrastructure types):
- Controllers, Custom Filters, Custom Middleware, Views, ViewModels, Startup

The `Startup` class / `Program.cs` is the **composition root** — it wires implementation types to interfaces, enabling dependency injection at runtime. (For DI wiring the UI project may need to reference Infrastructure, but actual references to Infrastructure types must be limited to the composition root.)

**ardalis/CleanArchitecture** template projects: `src/Clean.Architecture.Application` (use cases), `src/Clean.Architecture.Domain` (entities), `src/Clean.Architecture.Infrastructure` (EF Core, repositories), `src/Clean.Architecture.Web` (UI + composition root), plus `tests/`. The main branch targets .NET 9 (NuGet 10.x); installable via `dotnet new` templates `clean-arch` (full) and `min-clean` (minimal).

### Pattern 4: Testing

- **Unit test the Application Core in isolation** — it has no dependency on Infrastructure.
- **Integration test Infrastructure implementations** against real external dependencies (e.g., real DB).
- Because the UI layer has no direct dependency on Infrastructure types, implementations can be swapped for tests (fakes/mocks) or in response to changing requirements.
- Cockburn: automated function regression tests (e.g., FIT + mock DB) detect any leak of business logic into presentation/infrastructure layers.

## Layer Responsibilities

| Circle (inner → outer) | Responsibility | Typical contents | Sensitivity to external change |
|---|---|---|---|
| Entities | Enterprise-wide business rules | Business objects, data structures + functions | Lowest |
| Use Cases | Application-specific rules; orchestration of entities | Use case classes, output ports (interfaces) | Low (affected only by changes to app operation) |
| Interface Adapters | Data conversion between inner formats and external formats | Controllers, Presenters, Views, Models, repositories/gateways, all SQL | Medium |
| Frameworks & Drivers | Glue code to the outside world | Database, web framework, UI toolkit, device drivers | Highest |

## Best Practices

1. **Enforce the Dependency Rule mechanically**: source dependencies only point inward; inner circles never mention outer names or outer data formats.
2. **Cross boundaries with interfaces (DIP)**, never with concrete outer types — the use case calls an output port, the outer layer implements it.
3. **Pass simple data structures (DTOs, structs, arguments) across boundaries** — never entities, never database rows, never framework formats (e.g., RowStructure).
4. **Keep business rules framework-free**: use frameworks as tools at the periphery; don't let them constrain the core.
5. **Design for testability**: business rules must run without UI, DB, web server — then write unit tests against the core and integration tests against real infrastructure.
6. **Restrict SQL and persistence specifics to the adapter/infrastructure layer**, behind repository interfaces owned by the core.
7. **Write use cases at the application boundary**, specifying functions/events independent of external technology (Cockburn) — shorter, more stable, cheaper to maintain.
8. **Use dependency injection with a composition root** to wire edge implementations to core interfaces at runtime.
9. **Externalize the database**: treat it as a storage service behind an interface, not as the center of the application.
10. **Keep layering logical, not physical**: layers are logical separation; deployment tiers are a separate concern (an N-layer app can deploy on a single tier).
11. **Don't over-apply**: this architecture suits long-lived business applications with complex behavior; it is not appropriate for small websites (Palermo). Start monolithic and split only when justified.
12. **Add layers as needed** — the four circles are schematic; the Dependency Rule is the invariant.

## Common Pitfalls

1. **Business logic leaking into the UI/presentation layer** — the recurring failure that layered architectures promise to fix but rarely enforce; without a detection mechanism, the new layer is cluttered with business logic within a few years (Cockburn). Mitigation: automated function regression tests that fail when logic appears outside the core.
2. **Passing database rows or framework data formats across boundaries** — forces an inner circle to know an outer circle; violates the Dependency Rule (Martin's RowStructure example).
3. **Inner circles referencing outer names** — importing framework classes/types into domain or use case code; the #1 mechanical violation to catch in review.
4. **Database-centric design** ("database applications"): the database is external; if the app only makes sense with the DB present, it cannot be tested or evolved in isolation.
5. **UI instantiating or statically calling Infrastructure types** outside the composition root (Microsoft's rules for the UI layer).
6. **Transitive coupling in naive N-layer setups**: UI → BLL → DAL chains the business logic to data-access details; testing business logic requires a test database until dependency inversion is applied (MS Learn).
7. **Dogmatic rigidity**: treating exactly four circles as a law; the circles are schematic — the rule that matters is the Dependency Rule.
8. **Over-engineering small projects**: applying Clean Architecture where a simple website suffices adds complexity without benefit.
9. **Use cases coupled to external technologies**: use cases written with intimate knowledge of the technology outside each port become long, brittle, and expensive to maintain (Cockburn).
10. **Passing Entities across boundaries** to "cheat" and avoid writing DTOs — makes every layer depend on the domain's persistence shape.

## Version Notes

- **2005-09-04** — Hexagonal (Ports & Adapters) Architecture, Alistair Cockburn, HaT Technical Report 2005.02 (v0.9). Updated/republished on alistair.cockburn.us.
- **2008-07-29** — The Onion Architecture, part 1, Jeffrey Palermo (parts 2–4 followed).
- **2011** — Screaming Architecture, Robert C. Martin.
- **2012-08-13** — *The Clean Architecture* blog post, Robert C. Martin — canonical statement of the methodology.
- **2017** — *Clean Architecture: A Craftsman's Guide to Software Structure and Design* (book), Robert C. Martin — full treatment including crossing boundaries, the circles, and testing.
- Microsoft Learn e-book *Architect Modern Web Applications with ASP.NET Core and Azure* (chapter "Common web application architectures", updated 2023-03-07) and the eShopOnWeb reference application apply the style in .NET; ardalis/CleanArchitecture (18.4k stars, MIT) is a production-proven ASP.NET Core template.
- The name lineage: Hexagonal → Ports-and-Adapters → Onion → Clean; "Clean Architecture is just the latest in a series of names for the same loosely-coupled, dependency-inverted architecture."
