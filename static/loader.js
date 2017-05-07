var apiUrl = window.location.pathname + "api/";
var tagUrl = window.location.pathname + "tags/?main=true";
var searchValue = window.location.search;
var currentTag = extractTag("");
var autoReloadDelay = 5000;
var isReloading = false;
var reloadRef;
var reloadSince = 0;
var ajaxImg = new XMLHttpRequest();
ajaxImg.onreadystatechange = function() {
	if( this.readyState == 4 && this.status == 200 ) {
		setImage(this.responseText);
	}
}
var ajaxTags = new XMLHttpRequest();
ajaxTags.onreadystatechange = function() {
	if( this.readyState == 4 && this.status == 200 ) {
		updateTags( this.responseText );
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

function toggleSettings() {
	var settingsDiv = document.getElementById("settings");
	var expander = document.getElementById("settingsExpander");
	if( settingsDiv.className == "closed" ) {
		settingsDiv.className = "open";
		expander.value = "-";
	} else {
		settingsDiv.className = "closed";
		expander.value = "+";
	}
}

function toggleAutoReload() {
	if( isReloading ) {
		window.clearTimeout( reloadRef );
	}
	isReloading = !isReloading;
	if( isReloading ) {
		autoReload( autoReloadDelay );
	}
}

function autoReload( delay ) {
	if( isReloading ) {
		reloadSince = Date.now();
		reloadRef = window.setTimeout( function() {
			if( isReloading ) {
				reloadSince = 0;
				reloadImage();
				autoReload( autoReloadDelay );
			}
		}, delay );
	}
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
	ajaxImg.open("GET", apiUrl + searchValue, true)
		ajaxImg.send();
}

function setup() {
	setupTags();
	reloadImage();
	autoReload( autoReloadDelay );
	setCurrentTag();
}

function setReloadDelay( newValue ) {
	autoReloadDelay = newValue * 1000;
	document.getElementById("delayShower").innerText = newValue;
}

function setCurrentTag( newTag ) {
	if(newTag == "" || !newTag) {
		searchValue = "";
	} else {
		searchValue = "?tag=" + newTag;
	}
}

function extractTag() {
	if( currentTag == "" ) {
		return "";
	}
	var captured = /tag=([^&]+)/.exec( searchValue );
	if( captured == undefined || captured[1] == undefined ) {
		return ""
	}
	return captured[1]
}

function setupTags() {
	console.log("Getting main tags");
	ajaxTags.open("GET", tagUrl, true)
		ajaxTags.send();
}

function appendTag( tag ) {
	var isSelected = "";
	if( currentTag == tag ) {
		isSelected = 'checked="checked"'
	}
	this.innerHTML += '<input type="radio" name="tagRadio"' + isSelected + 'onchange="setCurrentTag(this.value)" value="' + tag + '">' + tag + '</input><br />'
}

function updateTags( mainTagsJson ) {
	var isSelected = "";
	if( currentTag == "" ) {
		isSelected = ' checked="checked" '
	}
	mainTags = JSON.parse( mainTagsJson );
	tagsDiv = document.getElementById( "tagsDiv" );
	tagsDiv.innerHTML = '<input type="radio" name="tagRadio"' + isSelected + 'onchange="setCurrentTag(this.value)" value="">All</input><br />';
	mainTags.forEach(appendTag, tagsDiv);
}
