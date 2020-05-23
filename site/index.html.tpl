<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>{{index .Environments "KIN_TITLE"}}</title>
    <link rel="stylesheet" href="css/style.css">
    <meta http-equiv="refresh" content="10" >
  </head>
  <body>
    <div id='title'>{{index .Environments "KIN_TITLE"}}</div>
    <div id='container'>
        {{ range $group, $containers := .Containers -}}
        <div class="group"> 
          <div class="groupheader">{{ $group }}</div>
          <div class="groupcontent">
            {{ range $containers -}}
            <a href="{{default "." .URL}}" class="item {{.Type}}">{{.Name}} <span class="tag tag-{{.Type}}">{{.Type}}</span></a>
            {{end -}}
        </div>
        </div>
        {{end -}}
        </div>
      </body>
</html>