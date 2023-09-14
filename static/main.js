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

	navigator.clipboard.writeText(link);
	alert("Link copied to clipboard " + link);
}

/**
 * reads roomkey from current url
 * @return {string} RoomKey
 */
function getRoomKeyFromUrl() {
	const url = new URL(window.location.href);
	return url.pathname.split("/").pop();
}

document.addEventListener("DOMContentLoaded", function () {
	const roomKeyInput = document.getElementById("RoomKey");

	if (!roomKeyInput) {
		return;
	}

	roomKeyInput.value = getRoomKeyFromUrl();
	if (roomKeyInput.value != "") {
		roomKeyInput.type = "hidden";
	}
});
