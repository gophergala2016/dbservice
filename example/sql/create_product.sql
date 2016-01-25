insert into products (name, price{{if .status}}, status{{end}}) values ({{.name | quote}}, {{.price}}{{if .status}}, {{.status | quote }}{{end}}) returning *
