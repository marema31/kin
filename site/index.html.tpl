<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>pkged site</title>
    <link rel="stylesheet" href="css/style.css">
  </head>
  <body>
    <ul>
        {{ range $group, $containers := . }}
        <li> {{ $group }}
          <ul>
            {{ range $containers}}
            <li><a href="{{.URL}}" class="{{.Type}}">{{.Name}}</li>
            {{end}}
          </ul>
        </>
        {{end}}
    </ul>
  </body>
</html>