/**
 * toggles element inside a form.
 * identified by classes `edit-show` and `edit-hide`
 * @param {string} formId
 * @param {boolean} visible
 */
function setEditMode(formId, visible) {
	const form = document.getElementById(formId)


	if (form) {

		const inputs = form.querySelectorAll(".edit-show")
		const rendered = form.querySelectorAll(".edit-hide")

		rendered.forEach((el) => {
			el.style.display = visible ? "none" : "block"
		})

		inputs.forEach((el) => {
			el.style.display = visible ? "block" : "none"
		})
	}
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
		return
	}


	roomKeyInput.value = getRoomKeyFromUrl();
	if (roomKeyInput.value != "") {
		roomKeyInput.type = "hidden";
	}
});
