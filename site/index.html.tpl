<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>pkged site</title>
    <link rel="stylesheet" href="css/style.css">
  </head>
  <body>
    <ul>
        {{ range . }}
        <li><a href="{{.URL}}">{{.Name}}</li>
        {{end}}
    </ul>
  </body>
</html>