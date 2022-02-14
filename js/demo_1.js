const leftButton = document.getElementById("left-button")
const rightButton = document.getElementById("right-button")

leftButton.addEventListener('touchend', function (event) {
    leftButton.style.border = "2px solid red"
    rightButton.style.border = "2px solid white"
    window.open("/", "_self")
}, false);

rightButton.addEventListener('touchend', function (event) {
    console.log("clicked")
    leftButton.style.border = "2px solid white"
    rightButton.style.border = "2px solid red"
    window.open("/demo_2", "_self")
}, false);

