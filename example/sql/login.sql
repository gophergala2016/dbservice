select id, name, email, json_build_object('user_id', id) as __jwt from users where email = lower({{.params.email | quote }}) AND password = crypt({{.params.password | quote }}, password)
