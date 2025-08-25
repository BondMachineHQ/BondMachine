package bondirect

import (
	"bytes"
	"fmt"
	"text/template"
)

const meshDotTemplate = `digraph BondirectMesh {
    rankdir=LR;
    node [shape=box, style=filled, fillcolor=lightblue];
    edge [color=darkblue];

    // Nodes
{{- range .Nodes}}
    subgraph cluster_{{.PeerId}} {
	style="filled, rounded" fillcolor=springgreen3 color=grey30;
	label="BM {{.PeerId}}";
	{{ $peerName := .Name }}
	// "{{.Name}}" [label="{{.Name}}"];

{{- range $.Edges}}
{{- if eq $peerName .NodeA }}
	subgraph cluster_{{$peerName}}_{{.Name}} {
	style="filled, rounded" fillcolor=aquamarine1 color=grey30;
	label="{{.Name}}";
	"{{.NodeA}}{{.Name}}{{.AtoBTransA}}" [label="{{.AtoBTransA}}"];
	"{{.NodeA}}{{.Name}}{{.BtoATransA}}" [label="{{.BtoATransA}}"];
	}
{{- end}}
{{- if eq $peerName .NodeB }}
	subgraph cluster_{{$peerName}}_{{.Name}} {
	style="filled, rounded" fillcolor=aquamarine1 color=grey30;
	label="{{.Name}}";
	"{{.NodeB}}{{.Name}}{{.AtoBTransB}}" [label="{{.AtoBTransB}}"];
	"{{.NodeB}}{{.Name}}{{.BtoATransB}}" [label="{{.BtoATransB}}"];
	}
{{- end}}
{{- end}}

    }
{{- end}}

{{- range $.Edges}}
	"{{.NodeA}}{{.Name}}{{.AtoBTransA}}" -> "{{.NodeB}}{{.Name}}{{.AtoBTransB}}";
	"{{.NodeB}}{{.Name}}{{.BtoATransB}}" -> "{{.NodeA}}{{.Name}}{{.BtoATransA}}";
{{- end}}
}
`

type dotTemplateData struct {
	Nodes []dotNode
	Edges []dotEdge
}

type dotNode struct {
	PeerId int
	Name   string
}

type dotEdge struct {
	NodeA      string
	NodeB      string
	Name       string
	AtoBTransA string
	AtoBTransB string
	BtoATransA string
	BtoATransB string
}

func EmitMeshDot(c *Config, mesh *Mesh) (string, error) {
	// Create template
	tmpl, err := template.New("meshDot").Parse(meshDotTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Prepare template data
	templateData := dotTemplateData{
		Nodes: make([]dotNode, 0, len(mesh.Nodes)),
	}

	for peerName, node := range mesh.Nodes {
		dotNode := dotNode{
			PeerId: int(node.PeerId),
			Name:   peerName,
		}

		templateData.Nodes = append(templateData.Nodes, dotNode)
	}

	for edgeName, edge := range mesh.Edges {
		dotEdge := dotEdge{
			NodeA:      edge.NodeA,
			NodeB:      edge.NodeB,
			Name:       edgeName,
			AtoBTransA: edge.FromAtoB.ATransceiver,
			AtoBTransB: edge.FromAtoB.BTransceiver,
			BtoATransA: edge.FromBtoA.ATransceiver,
			BtoATransB: edge.FromBtoA.BTransceiver,
		}

		templateData.Edges = append(templateData.Edges, dotEdge)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
