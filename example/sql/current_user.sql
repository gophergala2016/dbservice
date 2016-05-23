{{ if .jwt.user_id }}
select id, name, email from users where id={{.jwt.user_id | quote}}
{{ else }}
select false as logged_in
{{ end }}
