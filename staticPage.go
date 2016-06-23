package main

import (
	"fmt"
)

func rootPage() []byte {
	return []byte(fmt.Sprintf(rootTemplate))
}

const rootTemplate = `
<html>
<head>
<style type="text/css">
body {
	background: black;
	margin: 0px;
}

img {
	padding: 0px;
	margin: 0px;
}

img.scaledsize {
	height: 100%%;
}

img.originalsize {
}

</style>

<script type="text/javascript">
var ajax = new XMLHttpRequest();
var apiUrl = window.location.pathname + "api/" + window.location.search;
ajax.onreadystatechange = function() {
	if( this.readyState == 4 && this.status == 200 ) {
		setImage(this.responseText);
	}
}

function toggleSize() {
	var img = document.getElementById("randomImage")
	if( img.className == "originalsize" ) {
		img.className = "scaledsize";
	} else {
		img.className = "originalsize";
	}
}

function toggleReload() {
}

function setImage( randomImageMetaJson ) {
	var img = document.getElementById("randomImage");
	var loadingImg = new Image();
	img.onload = function() {
		img.src = this.src;
	}
	var randomImageMeta = JSON.parse( randomImageMetaJson );

	console.log("Setting image to " + randomImageMeta.url)
	img.alt = "source: " + randomImageMeta.source
	img.src = randomImageMeta.url;
}

function reloadImage() {
	console.log("reloading image");
	ajax.open("GET", apiUrl, true)
	ajax.send();
}

</script>
</head>
<body onload="reloadImage()">
<img id="randomImage" class="scaledsize" onclick="toggleSize" />
</body>
</html>
`
