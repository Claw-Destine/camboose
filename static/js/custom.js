function setActiveMenu(clickedElement) {
    const menuItems = document.querySelectorAll('#menu-list a');
    menuItems.forEach(item => item.classList.remove('is-active'));
    clickedElement.classList.add('is-active');
}
function showElement(elementId) {
    const element = document.querySelectorAll("#" + elementId);
    element.forEach(item => item.classList.add('is-active'));
}
function hideElement(elementId) {
    const element = document.querySelectorAll("#" + elementId);
    element.forEach(item => item.classList.remove('is-active'));
}
htmx.on("htmx:responseError", function (evt) {
    const requestConfig = evt.detail.requestConfig;
    const error = evt.detail.xhr.response;
    const msgElem = document.querySelectorAll("#notification-msg");
    msgElem.forEach(item => item.innerHTML = error);
    const notElem = document.querySelectorAll("#notification-bar");
    notElem.forEach(item => item.classList.remove('is-hidden'));
})
function hideNotification(elementId) {
    const element = document.querySelectorAll("#" + elementId);
    element.forEach(item => item.classList.add('is-hidden'));
}
function toggleClass(elementId, cls = "is-active") {
    function toggleItem(item) {
        if (item.classList.contains(cls)) {
            item.classList.remove(cls);
        } else {
            item.classList.add(cls);
        }
    }
    const element = document.querySelectorAll("#" + elementId);
    element.forEach(item => toggleItem(item));
}
