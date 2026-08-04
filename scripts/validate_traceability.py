#!/usr/bin/env python3
"""Validate traceability.ttl (OWL 2 DL, Turtle) — VEDO EduTrack.

The traceability model is the project's single source of truth for artifact
chains (vision → UC → US → FR/ADR → COMP → TEST). Per project rules it MUST
be touched only through Python tooling (rdflib), never by blind hand-editing.

Checks (exit 0 = pass, 1 = fail):
  1. Syntax: the file must parse as Turtle (rdflib raises on any syntax error).
  2. Structure: every instance must be of a known artifact class and carry a
     hasArtifactId.
  3. Orphan detection: every artifact instance (except roots: vision, and
     pure taxonomy classes) must participate in at least one traceability
     relation (refines/derivesFrom/partOf/isSourceOf/realizes/conformsTo/
     validates/relatesTo/checksPipeline/derivesFrom). Orphans are reported
     as failures — a chain must never dangle.
  4. File paths: hasFilePath values must point to existing files (relative to
     the repository root).

Usage (from repo root):
  uv run --with rdflib python scripts/validate_traceability.py [path/to/traceability.ttl]
"""
from __future__ import annotations

import os
import sys
from pathlib import Path

RD = "https://vedo-edutrack.dev/traceability#"

# Turtle predicates that count as "participating in a traceability chain".
LINK_PREDICATES = {
    RD + "refines",
    RD + "derivesFrom",
    RD + "partOf",
    RD + "isSourceOf",
    RD + "realizes",
    RD + "conformsTo",
    RD + "validates",
    RD + "relatesTo",
    RD + "checksPipeline",
    # Vision anchors (F1-F6): functions are linked when artifacts cite them.
    RD + "hasFunction",
    RD + "hasProblem",
    # Documentation links (DOC ↔ ADR/COMP).
    RD + "documents",
    RD + "isDocumentedBy",
}

# Root artifacts that are allowed to have no incoming/outgoing chain links.
ALLOWED_ROOTS = {"vision"}

# Classes that are taxonomy, not instances — never orphans.
TAXONOMY_CLASSES = {
    "Artifact", "Vision", "Glossary", "ArchitectureDecision", "UseCase",
    "UserStory", "Requirement", "FunctionalRequirement", "NonFunctionalRequirement",
    "Component", "Test", "UnitTest", "IntegrationTest", "E2EApiTest", "E2EGuiTest",
    "E2ETest", "LoadTest", "Gate", "Pipeline", "C4Diagram", "BusinessFunction", "BusinessProblem",
    "Documentation", "Interface", "DeploymentArtifact",
}


def main() -> int:
    try:
        import rdflib
        from rdflib import BNode, Graph, Namespace, RDF
    except ImportError as exc:
        print(f"[traceability] ERROR: rdflib is required: {exc}", file=sys.stderr)
        print("[traceability] install with: uv pip install rdflib", file=sys.stderr)
        return 1

    root = Path(__file__).resolve().parent.parent
    ttl_path = Path(sys.argv[1]) if len(sys.argv) > 1 else root / "traceability.ttl"
    ttl_path = ttl_path.resolve()

    if not ttl_path.exists():
        print(f"[traceability] ERROR: file not found: {ttl_path}", file=sys.stderr)
        return 1

    TR = Namespace(RD)
    g = Graph()

    # 1. Syntax — rdflib raises on invalid Turtle.
    try:
        g.parse(str(ttl_path), format="turtle")
    except Exception as exc:  # noqa: BLE001 — report any parse error and fail
        print(f"[traceability] FAIL: Turtle parse error: {exc}", file=sys.stderr)
        return 1

    triples = len(g)
    print(f"[traceability] OK: parsed {triples} triples from {ttl_path.name}")

    # 2+3. Structure + orphans.
    errors: list[str] = []
    # Anything declared as an OWL/RDF property is vocabulary, not an artifact.
    PROPERTY_TYPES = {
        RDF.Property,
        Namespace("http://www.w3.org/2002/07/owl#").ObjectProperty,
        Namespace("http://www.w3.org/2002/07/owl#").DatatypeProperty,
        Namespace("http://www.w3.org/2002/07/owl#").AnnotationProperty,
    }
    property_nodes = set()
    for pt in PROPERTY_TYPES:
        property_nodes |= set(g.subjects(RDF.type, pt))
    # Predicates (even undeclared) are vocabulary too.
    property_nodes |= {p for s, p, o in g}

    instance_nodes = set(g.subjects(RDF.type, None))
    # Blank nodes (e.g. owl:AllDisjointClasses members) are not artifacts.
    instance_nodes = {n for n in instance_nodes if not isinstance(n, BNode)}
    instance_nodes -= property_nodes

    for node in sorted(instance_nodes, key=str):
        local = str(node).split("#")[-1]
        # Skip the ontology node (empty local name), the namespace itself,
        # and pure taxonomy declarations.
        if local == "" or local in TAXONOMY_CLASSES or node == TR:
            continue
        # Every instance must carry a hasArtifactId (except the ontology + classes).
        if g.value(node, TR.hasArtifactId) is None:
            errors.append(f"{local}: missing tr:hasArtifactId")
            continue
        # Orphan check: must have at least one chain link (either direction).
        linked = False
        for p, _ in g.predicate_objects(node):
            if str(p) in LINK_PREDICATES:
                linked = True
                break
        if not linked:
            for _, p in g.subject_predicates(node):
                if str(p) in LINK_PREDICATES:
                    linked = True
                    break
        if not linked and local not in ALLOWED_ROOTS:
            errors.append(f"{local}: orphan artifact (no traceability link)")

    # 4. File paths (retired artifacts legitimately point to removed files).
    for node in instance_nodes:
        retired = g.value(node, TR.retired)
        if retired is not None and str(retired).lower() == "true":
            continue
        for path_val in g.objects(node, TR.hasFilePath):
            rel = str(path_val)
            if not (root / rel).exists():
                errors.append(f"{str(node).split('#')[-1]}: hasFilePath missing: {rel}")

    if errors:
        print(f"[traceability] FAIL: {len(errors)} issue(s):", file=sys.stderr)
        for err in errors:
            print(f"  - {err}", file=sys.stderr)
        return 1

    print("[traceability] OK: structure, orphans and file paths valid")
    return 0


if __name__ == "__main__":
    sys.exit(main())
