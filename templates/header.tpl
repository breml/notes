{{define "header"}}
<!doctype html>
<head>
 <meta charset="UTF-8">
 <title>notes</title>
 <link href="/_static/css/bootstrap.min.css" rel="stylesheet">
 <style>
  .nav>li>a { padding: 2px 3px; }
  textarea { font-family: Consolas, Lucida Console, monospace; }
  .page-header { margin-top: 0; margin-bottom: 0;}
  .page-header .buttons { padding-top: 31px; }
 </style>
 <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
 <div class="container">
  <div class="row col-md-12 page-header">
   <div class="col-md-3">
    <h1><a href="/" accesskey="h">Notes</a></h1>
   </div>
   <div class="col-md-9 buttons">
    <form method="POST">
     {{ if .Edit }}
      <button type="submit" class="btn btn-default btn-xs" name="edit" value="false" accesskey="c">
       <span class="glyphicon glyphicon-remove"></span> Cancel
      </button>
     {{ else }}
      <button type="submit" class="btn btn-default btn-xs" name="edit" value="true" accesskey="e">
       <span class="glyphicon glyphicon-edit"></span> Edit
      </button>
     {{ end }}
     <span>{{ .Url }}</span>
    </form>
   </div>
  </div>
  <div class="row col-md-12">
   <div class="col-md-3">
    <h4>Dirs</h4>
    <ol class="nav nav-pills nav-stacked">
     {{range $dir := .Dirs }}
      <li><a href="{{ $dir.Path }}">{{$dir.Name}}</a></li>
     {{ end }}
    </ol>
    <form method="POST">
     <div class="row">
      <div class="input-group-sm col-xs-9">
       <input type="text" class="form-control" name="file" placeholder="folder" accesskey="o"/>
      </div>
      <div class="input-group-sm col-xs-3">
       <button type="submit" class="btn btn-default" name="action" value="create_folder">
        <span class="glyphicon glyphicon-plus"></span>
       </button>
      </div>
     </div>
    </form>
    <h4>Files</h4>
    <ol class="nav nav-pills nav-stacked">
     {{range $file := .Files }}
      <li><a href="{{ $file.Path }}">{{$file.Name}}</a></li>
     {{ end }}
    </ol>
    <form method="POST">
     <div class="row">
      <div class="input-group-sm col-xs-9">
       <input type="text" class="form-control" name="file" placeholder="file" accesskey="f"/>
      </div>
      <div class="input-group-sm col-xs-3">
       <button type="submit" class="btn btn-default" name="action" value="create_file">
        <span class="glyphicon glyphicon-plus"></span>
       </button>
      </div>
     </div>
    </form>
   </div>
{{end}}