insert into users (name, email) values ({{.name | quote}}, {{.email | quote}}) returning *
