<!DOCTYPE html>

<html>
<head>
	<title>Beego</title>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
</head>

<body class="container">
<header class="text-center">
	<h1 class="logo">Welcome to Beego</h1>
	<div class="content" >
		{{range  $element := .posts}}
		<div style="margin-bottom: 40px; border-radius: 7px; background: lightgray">
			<div style="height: 400px;width: 100%;background-color: black; background-image: url({{$element.ImageUrl}}); background-size: contain;background-repeat: no-repeat;background-position: center;border-radius:
			 7px"></div>
		<div style="margin-bottom: 40px; padding:20px" class="text-left">
			<h3 style="margin-top: 5px">{{$element.Title}}</h3>
			<h5>by {{$element.Artist}}</h5>
			<p>link {{$element.Url}}</p>
			<p>#art #illustration #drawing #draw #picture #photography #artist #sketch #sketchbook #paper #pen #pencil #artsy #instaart #beautiful #instagood #gallery #masterpiece #creative #photooftheday #instaartist #graphic #graphics #artoftheday</p>
		</div>
		</div>
		{{end}}
	</div>
	{{if gt .paginator.PageNums 1}}
	<ul class="pagination pagination-sm">
		{{if .paginator.HasPrev}}
		<li><a href="{{.paginator.PageLinkFirst}}">First</a></li>
		<li><a href="{{.paginator.PageLinkPrev}}">&lt;</a></li>
		{{else}}
		<li class="disabled"><a>Fist</a></li>
		<li class="disabled"><a>&lt;</a></li>
		{{end}}
		{{range $index, $page := .paginator.Pages}}
		<li
				{{if $.paginator.IsActive .}} class="active" {{end}}>
			<a href="{{$.paginator.PageLink $page}}">{{$page}}</a>
		</li>
		{{end}}
		{{if .paginator.HasNext}}
		<li><a href="{{.paginator.PageLinkNext}}">&gt;</a></li>
		<li><a href="{{.paginator.PageLinkLast}}">Last</a></li>
		{{else}}
		<li class="disabled"><a>&gt;</a></li>
		<li class="disabled"><a>Last</a></li>
		{{end}}
	</ul>
	{{end}}
</header>

<div class="backdrop"></div>

<script src="/static/js/reload.min.js"></script>
</body>
</html>
