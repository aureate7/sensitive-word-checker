package main

import "sort"

type FuzzyMatcher struct {
	normalizer *Normalizer
}

func NewFuzzyMatcher(normalizer *Normalizer) *FuzzyMatcher {
	return &FuzzyMatcher{normalizer: normalizer}
}

type AliasIndex struct {
	Automaton *ACAutomaton
	AliasToKW map[string][]string
}

func (m *FuzzyMatcher) BuildAliasIndex(words map[string]struct{}) *AliasIndex {
	aliasToKW := make(map[string][]string)
	for w := range words {
		alias := m.normalizer.NormalizeTextAggressive(w).NormalizedText
		if alias == "" {
			continue
		}
		aliasToKW[alias] = appendUniq(aliasToKW[alias], w)
	}
	aliases := make([]string, 0, len(aliasToKW))
	for a := range aliasToKW {
		aliases = append(aliases, a)
	}
	sort.Strings(aliases)
	ac := NewAC()
	ac.Build(aliases)
	return &AliasIndex{Automaton: ac, AliasToKW: aliasToKW}
}

func appendUniq(in []string, v string) []string {
	for _, x := range in {
		if x == v {
			return in
		}
	}
	return append(in, v)
}
