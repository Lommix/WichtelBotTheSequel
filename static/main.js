/**
 * toggles element inside a form.
 * identified by classes `edit-show` and `edit-hide`
 * @param {Event} event
 * @param {string} parentId
 * @param {boolean} visible
 */
function setEditMode(event, parentId, visible) {
	event.preventDefault();
	const form = document.getElementById(parentId);
	if (form) {
		const inputs = form.querySelectorAll(".edit-show");
		const rendered = form.querySelectorAll(".edit-hide");

		rendered.forEach((el) => {
			el.style.display = visible ? "none" : "block";
		});

		inputs.forEach((el) => {
			el.style.display = visible ? "block" : "none";
		});
	}
}

/**
 * toggles element visibility
 * @param {Event} event
 * @param {string} element_id
 * @param {boolean} state
 */
function setVisibilty(event, id, state) {
	event.preventDefault();
	const element = document.getElementById(id);
	if (element) {
		element.style.display = state ? "block" : "none";
	}
}

/**
 * saves key to local storage
 * @param {string} key
 */
function saveRoomKey(key) {
	localStorage.setItem("room", key);
}

/**
 * saves key to local storage
 * @return {string} key
 */
function loadRoomKey() {
	return localStorage.getItem("room") || "";
}

/**
 * toggles element visibility
 * @param {Event} event
 * @param {string} element_id
 */
function toggleVisibilty(event, id) {
	event.preventDefault();
	const element = document.getElementById(id);
	if (element) {
		element.style.display =
			element.style.display === "block" ? "none" : "block";
	}
}

/**
 * copies link to clipboard
 * @param {Event} event
 * @param {DOMElement} element
 */
function saveToClipboard(event, element) {
	event.preventDefault();
	const link = element.href;
	//desktop
	navigator.clipboard.writeText(link);
	//mobile
	const clipElement = document.createElement("input");
	clipElement.value = link;
	clipElement.classList.value = "p-1 absolute w-full h-fit left-0 top-0";
	clipElement.onblur = () => {
		element.removeChild(clipElement)
	}
	element.appendChild(clipElement);
	clipElement.select();
}

/**
 * reads roomkey from current url
 * @return {string} RoomKey
 */
function getRoomKeyFromUrl() {
	const url = new URL(window.location.href);

	if (url.toString().includes("/join/")) {
		return url.pathname.split("/join/").pop();
	}

	if (url.toString().includes("/login/")) {
		return url.pathname.split("/login/").pop();
	}

	return "";
}

document.addEventListener("DOMContentLoaded", function () {
	const roomNode = document.getElementById("RoomContainer");
	if (roomNode) {
		saveRoomKey(roomNode.dataset.key);
	}

	const roomKeyInput = document.getElementById("RoomKey");
	if (!roomKeyInput) {
		return;
	}

	roomKeyInput.value = loadRoomKey();
	urlKey = getRoomKeyFromUrl();
	if (urlKey) {
		console.log(urlKey);
		roomKeyInput.value = urlKey;
	}

	if (roomKeyInput.value != "") {
		roomKeyInput.type = "hidden";
	}
});
