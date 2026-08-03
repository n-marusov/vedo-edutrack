package mockhub

import (
	"strings"
	"testing"
)

func TestParseTraceabilityTTL(t *testing.T) {
	// Mini-Turtle with the supported subset (mirrors traceability.ttl shape).
	src := `
@prefix owl: <http://www.w3.org/2002/07/owl#> .
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
@prefix tr: <https://vedo-edutrack.dev/traceability#> .

tr:Artifact a owl:Class ;
  rdfs:label "Artifact"@en ;
  rdfs:comment "Top-level class."@en .

tr:Vision a owl:Class ;
  rdfs:label "Vision"@en ;
  rdfs:subClassOf tr:Artifact .

tr:hasRequirement a owl:ObjectProperty ;
  rdfs:label "has requirement"@en ;
  rdfs:domain tr:Vision ;
  rdfs:range tr:Artifact ;
  owl:FunctionalProperty true .
`
	ont, err := Parse(strings.NewReader(src))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	classes, props := ont.Counts()
	if classes != 2 {
		t.Errorf("classes = %d, want 2", classes)
	}
	if props != 1 {
		t.Errorf("properties = %d, want 1", props)
	}

	artifact := ont.Classes["https://vedo-edutrack.dev/traceability#Artifact"]
	if artifact == nil {
		t.Fatal("Artifact class missing")
	}
	if artifact.Label != "Artifact" {
		t.Errorf("Artifact label = %q, want Artifact", artifact.Label)
	}

	vision := ont.Classes["https://vedo-edutrack.dev/traceability#Vision"]
	if vision == nil {
		t.Fatal("Vision class missing")
	}
	if len(vision.Parents) != 1 || vision.Parents[0] != artifact.ID {
		t.Errorf("Vision parents = %v, want [Artifact]", vision.Parents)
	}
	if len(artifact.Children) != 1 || artifact.Children[0] != vision.ID {
		t.Errorf("Artifact children = %v, want [Vision]", artifact.Children)
	}

	prop := ont.Properties["https://vedo-edutrack.dev/traceability#hasRequirement"]
	if prop == nil {
		t.Fatal("hasRequirement property missing")
	}
	if !prop.Functional {
		t.Error("hasRequirement should be functional")
	}
	if len(prop.Domains) != 1 || prop.Domains[0] != vision.ID {
		t.Errorf("hasRequirement domains = %v, want [Vision]", prop.Domains)
	}

	// Breadcrumb from Vision → Artifact.
	bc := ont.Breadcrumb(vision.ID)
	if len(bc) != 2 || bc[0] != "Artifact" || bc[1] != "Vision" {
		t.Errorf("breadcrumb = %v, want [Artifact Vision]", bc)
	}

	// Descendants of Artifact → Vision.
	desc := ont.Descendants(artifact.ID, 0)
	if len(desc) != 1 || desc[0] != vision.ID {
		t.Errorf("descendants = %v, want [Vision]", desc)
	}

	// Autocomplete.
	ac := ont.Autocomplete("artifact", 10)
	if len(ac) == 0 {
		t.Error("autocomplete('artifact') returned nothing")
	}
}

func TestParseErrorLine(t *testing.T) {
	src := "@prefix owl: <http://www.w3.org/2002/07/owl#> .\n<http://x> a .\n"
	_, err := Parse(strings.NewReader(src))
	if err == nil {
		t.Fatal("expected parse error")
	}
	if !strings.Contains(err.Error(), "line 2") {
		t.Errorf("error should carry line number, got %q", err)
	}
}
