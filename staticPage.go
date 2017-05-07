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
	<link rel="stylesheet" href="/static/style.css">
	<script type="text/javascript" src="/static/loader.js"></script>
</head>
<body onload="setup()">
<img id="randomImage" class="scaledsize" onclick="toggleSize()" />
<div id="settings" class="open">
	<input id="settingsExpander" type="button" onclick="toggleSettings()" value="-" /><br />
	<div>
		<div>
		<label><small>Auto reload</small><input type="checkbox" onchange="toggleAutoReload()"></input></label>
		</div>
		<div id=delayDiv>
			<small>Reload time</small>
			<span id="delaySpan">
				<input id="delaySlider" type="range" oninput="setReloadDelay(this.valueAsNumber)" min="1" max="30" value="5" step=1 />
				<span id="delayShower">5</span>
			</span>
		</div>
	</div>
	<div>
		<input type="button" onclick="reloadImage()" value="New image" />
	</div>
	<div id="tagsDiv">
	</div>
</div>
</body>
</html>
`
