package gapcoverage

import (
	"math/rand"
	"testing"
)

// BenchmarkDiagnoseRootCause1kModules measures root-cause diagnosis on a
// 1k-module graph (NFR-FR-execute.gap.diagnose-root-cause AC-4: <=2s).
func BenchmarkDiagnoseRootCause1kModules(b *testing.B) {
	graph := buildGapGraph(1000, 42)
	mastery := gapMastery(1000, 0.3)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DiagnoseRootCause(graph, mastery, "m999")
	}
}

func buildGapGraph(n int, seed int64) Graph {
	r := rand.New(rand.NewSource(seed))
	modules := make([]Module, 0, n)
	links := make([]Link, 0, n*2)
	for i := 0; i < n; i++ {
		modules = append(modules, Module{ID: "m" + itoa(i)})
		if i > 0 {
			links = append(links, Link{SourceID: "m" + itoa(i-1), TargetID: "m" + itoa(i), Type: LinkStrictPrerequisite})
		}
	}
	for i := 0; i < n; i++ {
		for j := i + 2; j < n && j < i+5; j++ {
			if r.Intn(3) == 0 {
				links = append(links, Link{SourceID: "m" + itoa(i), TargetID: "m" + itoa(j), Type: LinkStrictPrerequisite})
			}
		}
	}
	return Graph{Modules: modules, Links: links}
}

func gapMastery(n int, masteredRatio float64) Mastery {
	modules := map[string]float64{}
	for i := 0; i < n; i++ {
		if float64(i%100)/100 < masteredRatio {
			modules["m"+itoa(i)] = 1.0
		} else {
			modules["m"+itoa(i)] = 0.2
		}
	}
	return Mastery{Modules: modules}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(buf[pos:])
}
