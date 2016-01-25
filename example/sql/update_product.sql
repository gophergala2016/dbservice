update products set name={{.name | quote}}, price={{.price}}{{ if .status}}, status={{.status | quote }}{{end}} where id={{.id}} returning *
