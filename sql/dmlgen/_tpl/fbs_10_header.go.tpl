// Auto generated via github.com/weiwolves/pkg/sql/dmlgen

namespace {{.Package}};

include "github.com/weiwolves/pkg/storage/null/null.fbs";

{{range $opts := .SerializerHeaderOptions -}}
{{$opts}};
{{end}}
