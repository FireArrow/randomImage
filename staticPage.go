package main

import (
	"fmt"
)

func rootPage(img string, source string) []byte {
	return []byte(fmt.Sprintf(rootTemplate, img, source))
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
function toggleSize() {
	var img = document.getElementById("randomImage")
	if( img.className == "originalsize" ) {
		img.className = "scaledsize";
	} else {
		img.className = "originalsize";
	}
}
</script>
</head>
<body>
<img id="randomImage" class="scaledsize" onclick="toggleSize()" src="%s" alt="source: %s" />
</body>
</html>
`
