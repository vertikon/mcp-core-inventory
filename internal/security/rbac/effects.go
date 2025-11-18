package rbac

// PolicyEffect describes the result of a policy/permission evaluation.
type PolicyEffect string

const (
	// EffectAllow grants access when the conditions match.
	EffectAllow PolicyEffect = "allow"
	// EffectDeny blocks access when the conditions match.
	EffectDeny PolicyEffect = "deny"
)
