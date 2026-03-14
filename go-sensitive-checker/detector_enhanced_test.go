package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizerAggressive(t *testing.T) {
	n := NewNormalizer()
	res := n.NormalizeTextAggressive("你这个 s*b 真离谱")
	if res.NormalizedText != "你这个sb真离谱" {
		t.Fatalf("unexpected normalized text: %q", res.NormalizedText)
	}
}

func TestTermMappingToggle(t *testing.T) {
	n := NewNormalizer()

	enabled := DetectOptions{
		ExactMatch:        true,
		NormalizeMatch:    true,
		FuzzyMatch:        true,
		PinyinMatch:       true,
		EnableTermMapping: true,
		MappingMode:       MappingModeIncremental,
	}
	disabled := enabled
	disabled.EnableTermMapping = false

	resEnabled := n.NormalizeTextAggressiveWithOptions("你这个 s@b 真离谱", enabled)
	if resEnabled.NormalizedText != "你这个sab真离谱" {
		t.Fatalf("mapping enabled unexpected: %q", resEnabled.NormalizedText)
	}

	resDisabled := n.NormalizeTextAggressiveWithOptions("你这个 s@b 真离谱", disabled)
	if resDisabled.NormalizedText != "你这个sb真离谱" {
		t.Fatalf("mapping disabled unexpected: %q", resDisabled.NormalizedText)
	}
}

func TestCustomTermMappingOverride(t *testing.T) {
	n := NewNormalizer()
	opts := DetectOptions{
		ExactMatch:        true,
		NormalizeMatch:    true,
		FuzzyMatch:        true,
		PinyinMatch:       true,
		EnableTermMapping: true,
		MappingMode:       MappingModeOverride,
		CustomMappings: []TermMapping{
			{From: "@", To: "o"},
			{From: "vv", To: "w"},
		},
	}

	res := n.NormalizeTextAggressiveWithOptions("vva@b", opts)
	if res.NormalizedText != "waob" {
		t.Fatalf("unexpected override mapping text: %q", res.NormalizedText)
	}
}

func TestDetectPinyinAndExplain(t *testing.T) {
	base := setupTestWordRepo(t)
	cfg := DefaultDetectorConfig(base)
	cfg.PinyinAliasPath = filepath.Join(base, "拼音混淆词", "拼音映射.txt")
	d := NewDetectorWithConfig(base, cfg)

	resp := d.DetectWithOptions("你真是sha_bi", []string{AbusiveHigh}, &DetectOptions{
		ExactMatch:     true,
		NormalizeMatch: true,
		FuzzyMatch:     true,
		PinyinMatch:    true,
	})

	if !resp.HasSensitive {
		t.Fatalf("expected sensitive hit, got none")
	}
	if resp.NormalizedText != "你真是shabi" {
		t.Fatalf("unexpected normalized text: %q", resp.NormalizedText)
	}

	hit := findWord(resp.DetectedWords, "傻逼")
	if hit == nil {
		t.Fatalf("expected hit word 傻逼, got %+v", resp.DetectedWords)
	}
	if hit.CountPinyin == 0 {
		t.Fatalf("expected pinyin count > 0, got %+v", *hit)
	}
	if !containsStr(hit.MatchMethods, matchTypePinyin) {
		t.Fatalf("expected match type %s in %+v", matchTypePinyin, hit.MatchMethods)
	}

	foundEvidence := false
	for _, e := range resp.HitEvidences {
		if e.Word == "傻逼" && e.MatchType == matchTypePinyin {
			foundEvidence = true
			break
		}
	}
	if !foundEvidence {
		t.Fatalf("expected pinyin evidence in %+v", resp.HitEvidences)
	}
}

func TestDetectCanDisablePinyin(t *testing.T) {
	base := setupTestWordRepo(t)
	cfg := DefaultDetectorConfig(base)
	cfg.PinyinAliasPath = filepath.Join(base, "拼音混淆词", "拼音映射.txt")
	d := NewDetectorWithConfig(base, cfg)

	resp := d.DetectWithOptions("你真是sha_bi", []string{AbusiveHigh}, &DetectOptions{
		ExactMatch:     true,
		NormalizeMatch: true,
		FuzzyMatch:     true,
		PinyinMatch:    false,
	})

	if resp.HasSensitive {
		t.Fatalf("expected no hit when pinyin disabled, got %+v", resp)
	}
}

func TestDetectAutoPinyinForAllLoadedWords(t *testing.T) {
	base := t.TempDir()
	mustWriteFile(t, filepath.Join(base, "辱骂类敏感词", "辱骂高敏感词（添加版）.txt"), "混蛋\n")

	cfg := DefaultDetectorConfig(base)
	cfg.PinyinAliasPath = "" // 仅验证自动拼音生成，不依赖手工映射
	d := NewDetectorWithConfig(base, cfg)

	resp := d.DetectWithOptions("你这个 hun-dan 真离谱", []string{AbusiveHigh}, &DetectOptions{
		ExactMatch:     true,
		NormalizeMatch: true,
		FuzzyMatch:     true,
		PinyinMatch:    true,
	})
	if !resp.HasSensitive {
		t.Fatalf("expected auto pinyin hit, got %+v", resp)
	}

	hit := findWord(resp.DetectedWords, "混蛋")
	if hit == nil {
		t.Fatalf("expected hit word 混蛋, got %+v", resp.DetectedWords)
	}
	if hit.CountPinyin == 0 {
		t.Fatalf("expected pinyin count > 0, got %+v", *hit)
	}
}

func setupTestWordRepo(t *testing.T) string {
	t.Helper()
	base := t.TempDir()
	mustWriteFile(t, filepath.Join(base, "辱骂类敏感词", "辱骂高敏感词（添加版）.txt"), "傻逼\nsb\n")
	mustWriteFile(t, filepath.Join(base, "拼音混淆词", "拼音映射.txt"), "# word=alias1|alias2\n傻逼=shabi|sha bi|sha_bi|sha-bi\n")
	return base
}

func mustWriteFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file failed: %v", err)
	}
}

func findWord(words []WordHit, w string) *WordHit {
	for i := range words {
		if words[i].Word == w {
			return &words[i]
		}
	}
	return nil
}

func containsStr(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
