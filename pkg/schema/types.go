package schema

type Object struct {
	Name       string     `json:"name" yaml:"name"`
	Properties []string   `json:"properties" yaml:"properties"`
	UniqueOn   [][]string `json:"unique_on" yaml:"unique_on"`
}

type RelationshipType string

const OneToOne = RelationshipType("one-to-one")
const OneToMany = RelationshipType("one-to-many")
const ManyToOne = RelationshipType("many-to-one")
const ManyToMany = RelationshipType("many-to-many")

type Relationship struct {
	Name        *string          `json:"name" yaml:"name"`
	ReverseName *string          `json:"reverse_name" yaml:"reverse_name"`
	Source      string           `json:"source" yaml:"source"`
	Destination string           `json:"destination" yaml:"destination"`
	Type        RelationshipType `json:"type" yaml:"type"`
	Optional    bool             `json:"optional" yaml:"optional"`
}

type Schema struct {
	Objects       []*Object       `json:"objects" yaml:"objects"`
	Relationships []*Relationship `json:"relationships" yaml:"relationships"`
}
