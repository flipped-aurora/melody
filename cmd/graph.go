package cmd

import (
	"bytes"
	"github.com/spf13/cobra"
	"io"
	"melody/config"
	"os"
	"text/template"
)

func graphFunc(cmd *cobra.Command, args []string) {
	if cfgFilePath == "" {
		cmd.Println("Please provide the path to your melody config file ")
		return
	}
	serviceConfig, err := parser.Parse(cfgFilePath)
	if err != nil {
		cmd.Printf("ERROR parsing the melody config file: %s\n", err.Error())
		os.Exit(-1)
	}
	writeToDot(os.Stdout, serviceConfig, cmd)
}

func writeToDot(writer io.Writer, config config.ServiceConfig, cmd *cobra.Command) {
	t := template.New("dot")
	var buf bytes.Buffer
	if err := template.Must(t.Parse(tmplGraph)).Execute(&buf, config); err != nil {
		cmd.Println("ERROR convert data to dot error:", err)
	}
	buf.WriteTo(writer)
}

const tmplGraph = `digraph melody { {{ $port := .Port }}
    label="Melody Gateway";
    labeljust="l";
    fontname="Ubuntu";
    fontsize="13";
    rankdir="LR";
    bgcolor="aliceblue";
    style="solid";
    penwidth="0.5";
    pad="0.0";
    nodesep="0.35";

    node [shape="ellipse" style="filled" fillcolor="honeydew" fontname="Ubuntu" penwidth="1.0" margin="0.05,0.0"];

	{{ range $i, $endpoint := .Endpoints }}
    {{printf "subgraph \"cluster_%s\" {" .Endpoint }}
    	label="{{ .Endpoint }}";
    	bgcolor="lightgray";
    	shape="box";
    	style="solid";

        "{{ .Endpoint }}" [ shape=record, label="{ { Timeout | {{.Timeout.String}} } | { CacheTTL | {{.CacheTTL.String}} } | { Output | {{.OutputEncoding}} } | { QueryString | {{.QueryString}} } }" ]
        {{ if .ExtraConfig }}"extra_{{$i}}" [ shape=record, label="{ {ExtraConfig} {{ range $key, $value := .ExtraConfig }} | { {{ $key }} {{ range $k, $v := $value }}| { {{$k}} | {{$v}} } {{ end }} }{{ end }} }" ]{{ end }}
	    {{ range $j, $backend := .Backends }}
	    {{printf "subgraph \"cluster_%s\" {" .URLPattern }}
	    	label="{{ .URLPattern }}";
	    	bgcolor="beige";
	    	shape="box";
	    	style="solid";
        	"in_{{$i}}_{{$j}}" [ shape=record, label="{ {sd|{{ if .SD }}{{ .SD }}{{ else }}static{{ end }} } | { Hosts | {{.Host}} } | { Encoding | {{ if .Encoding }}{{ .Encoding }}{{ else }}JSON{{ end }} } }" ]
        {{ if .ExtraConfig }}"extra_{{$i}}_{{$j}}" [ shape=record, label="{ { ExtraConfig {{ range $key, $v := .ExtraConfig }} | {{ $key }} {{ end }} } }" ]{{ end }}
	    {{println "}" }}
	    "{{ $endpoint.Endpoint }}" -> in_{{$i}}_{{$j}} [ label="x{{ .ConcurrentCalls }}"]{{ end }}
    {{ println "}" }}{{ end }}
    {{ range .Endpoints }}
    ":{{ $port }}" -> "{{ .Endpoint }}" [ label="{{ .Method }}"]{{ end }}
}
`

