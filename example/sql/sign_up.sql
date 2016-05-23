INSERT INTO users (name, email, password) VALUES
  ({{.params.name | quote}}, {{.params.email | quote}}, crypt({{.params.password | quote}}, gen_salt('bf', 10)));
