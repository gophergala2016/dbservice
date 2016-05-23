insert into products (name, price{{if .params.status}}, status{{end}}) values ({{.params.name | quote}}, {{.params.price}}{{if .params.status}}, {{.params.status | quote }}{{end}}) returning *
