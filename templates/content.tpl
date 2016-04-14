{{ template "header" . }}
     <div class="col-md-9">
     {{ if .Edit }}
      <form method="POST">
       <div class="form-group">
        <textarea type="text" class="form-control" rows="15" placeholder="Insert markdown here" name="content" autofocus>{{ .Content }}</textarea>
       </div>
       <div class="form-group">
        <button type="submit" class="btn btn-default" name="action" value="save">
         <span class="glyphicon glyphicon-floppy-disk" accesskey="s"></span> Save
        </button>
        <button type="submit" class="btn btn-default" name="action" value="save_edit">
         <span class="glyphicon glyphicon-floppy-disk" accesskey="e"></span> Save & Edit
        </button>
       </div>
      </form>
     {{ end }}
      <div class="row col-md-12">
       {{ .ContentHTML }}
      </div>
     </div>
{{ template "footer" . }}