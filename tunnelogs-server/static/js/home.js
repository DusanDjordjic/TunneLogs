const input = document.getElementById("lobby-input");
const button = document.getElementById("sumbit-button");

button.addEventListener("click", function(event) {
    event.preventDefault();

    if (input.value == "") {
        console.error("Empty lobby name");
        return;
    }

    console.log("Navigating...", input.value);
    window.location.href = "/lobby/" + input.value;
});
