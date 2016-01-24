update users set name={{.name | quote}}, email={{.email | quote}} where id={{.id}}
