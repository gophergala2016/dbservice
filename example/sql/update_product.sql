update products set name={{.params.name | quote}}, price={{.params.price}}{{ if .params.status}}, status={{.params.status | quote }}{{end}} where id={{.params.id}} returning *
