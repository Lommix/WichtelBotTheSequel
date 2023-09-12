/**
 * @return {string} RoomKey
 */
function getRoomKeyFromUrl() {
	const url = new URL(window.location.href);
	return url.pathname.split("/").pop();
}

document.addEventListener("DOMContentLoaded", function () {
	const roomKeyInput = document.getElementById("RoomKey");
	console.log(roomKeyInput);
	console.log(getRoomKeyFromUrl());

	roomKeyInput.value = getRoomKeyFromUrl();
	if (roomKeyInput.value != "") {
		roomKeyInput.type = "hidden";
	}
});
