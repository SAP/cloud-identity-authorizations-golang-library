package dcn

type DcnContainer struct {
	Version   int        `json:"version"`
	Policies  []Policy   `json:"policies"`
	Functions []Function `json:"functions"`
	Schemas   []Schema   `json:"schemas"`
	Tests     []Test     `json:"tests"`
}

type Policy struct {
	QualifiedName QualifiedName  `json:"policy"`
	Rules         []Rule         `json:"rules"`
	Uses          []Use          `json:"uses"`
	Description   string         `json:"description,omitempty"`
	Default       bool           `json:"default,omitempty"`
	Internal      bool           `json:"internal,omitempty"`
	Annotations   map[string]any `json:"annotations,omitempty"`
}

type Rule struct {
	Effect      string         `json:"effect"`
	Actions     []string       `json:"actions"`
	Resources   []string       `json:"resources"`
	Condition   *Expression    `json:"condition,omitempty"`
	Role        bool           `json:"role,omitempty"`
	Annotations map[string]any `json:"annotations,omitempty"`
}

type Use struct {
	QualifiedPolicyName QualifiedName  `json:"use"`
	Restrictions        [][]Expression `json:"restrictions"`
	Annotations         map[string]any `json:"annotations,omitempty"`
}

type Function struct {
	QualifiedName QualifiedName  `json:"function"`
	Result        Expression     `json:"result"`
	ReturnType    string         `json:"returnType"`
	Internal      bool           `json:"internal,omitempty"`
	Annotations   map[string]any `json:"annotations,omitempty"`
}

type Schema struct {
	QualifiedName QualifiedName   `json:"schema"`
	Tenant        string          `json:"tenant,omitempty"`
	Definition    SchemaAttribute `json:"definition,omitempty"`
	Annotations   map[string]any  `json:"annotations,omitempty"`
}

type SchemaAttribute struct {
	Type        string                     `json:"attribute"`
	Nested      map[string]SchemaAttribute `json:"nested,omitempty"`
	Annotations map[string]any             `json:"annotations,omitempty"`
}

type QualifiedName []string

type Expression struct {
	Constant any

	Ref []string `json:"ref,omitempty"`

	Call []string     `json:"call,omitempty"`
	Args []Expression `json:"args,omitempty"`
}

type Test struct {
	Test       []string    `json:"test"`
	Assertions []Assertion `json:"assertions"`
}

type Assertion struct {
	Expect      Expression `json:"expect"`
	Policies    [][]string `json:"policies"`
	ScopeFilter [][]string `json:"scopeFilter"`
	Inputs      []Input    `json:"inputs"`
	Actions     []string   `json:"actions,omitempty"`
	Resources   []string   `json:"resources,omitempty"`
}

type Input struct {
	Input    map[string]any `json:"input"`
	Unknowns []Reference    `json:"unknowns"`
	Ignores  []Reference    `json:"ignores"`
}

type Reference struct {
	Ref []string `json:"ref"`
}
