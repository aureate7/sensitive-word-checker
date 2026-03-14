package main

import "os"

const (
	matchTypeExactRaw        = "exact_raw"
	matchTypeExactNoSymbol   = "exact_no_symbol"
	matchTypeExactNormalized = "exact_normalized"
	matchTypeFuzzy           = "fuzzy"
	matchTypePinyin          = "pinyin"

	MappingModeIncremental = "incremental"
	MappingModeOverride    = "override"
)

type TermMapping struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// DetectOptions controls optional strategies in detection pipeline.
type DetectOptions struct {
	ExactMatch        bool          `json:"exact_match"`
	NormalizeMatch    bool          `json:"normalize_match"`
	FuzzyMatch        bool          `json:"fuzzy_match"`
	PinyinMatch       bool          `json:"pinyin_match"`
	EnableTermMapping bool          `json:"enable_term_mapping"`
	MappingMode       string        `json:"mapping_mode,omitempty"` // incremental / override
	CustomMappings    []TermMapping `json:"custom_mappings,omitempty"`
}

func DefaultDetectOptions() DetectOptions {
	return DetectOptions{
		ExactMatch:        true,
		NormalizeMatch:    true,
		FuzzyMatch:        true,
		PinyinMatch:       true,
		EnableTermMapping: true,
		MappingMode:       MappingModeIncremental,
	}
}

func (o DetectOptions) FillDefault(def DetectOptions) DetectOptions {
	if !o.ExactMatch && !o.NormalizeMatch && !o.FuzzyMatch && !o.PinyinMatch {
		return def
	}
	return o
}

type DetectorConfig struct {
	DefaultOptions       DetectOptions
	EnableNormalize      bool
	EnableFuzzy          bool
	EnablePinyin         bool
	EnableAutoPinyin     bool
	EnablePinyinInitials bool
	PinyinAliasPath      string
}

func DefaultDetectorConfig(basePath string) DetectorConfig {
	return DetectorConfig{
		DefaultOptions:       DefaultDetectOptions(),
		EnableNormalize:      envBool("SENSITIVE_ENABLE_NORMALIZE", true),
		EnableFuzzy:          envBool("SENSITIVE_ENABLE_FUZZY", true),
		EnablePinyin:         envBool("SENSITIVE_ENABLE_PINYIN", true),
		EnableAutoPinyin:     envBool("SENSITIVE_ENABLE_AUTO_PINYIN", true),
		EnablePinyinInitials: envBool("SENSITIVE_ENABLE_PINYIN_INITIALS", false),
		PinyinAliasPath:      envStr("SENSITIVE_PINYIN_ALIAS_FILE", basePath+"/拼音混淆词/拼音映射.txt"),
	}
}

func envBool(key string, def bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	switch v {
	case "1", "true", "TRUE", "yes", "YES", "on", "ON":
		return true
	case "0", "false", "FALSE", "no", "NO", "off", "OFF":
		return false
	default:
		return def
	}
}

func envStr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}
