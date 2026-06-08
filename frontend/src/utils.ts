function setActiveMenu(clickedElement: Element) {
    const menuItems = document.querySelectorAll('#menu-list a');
    menuItems.forEach(item => item.classList.remove('is-active'));
    clickedElement.classList.add('is-active');
}
function showElement(elementId: string) {
    const element = document.querySelectorAll("#" + elementId);
    element.forEach(item => item.classList.add('is-active'));
}
function hideElement(elementId: string) {
    const element = document.querySelectorAll("#" + elementId);
    element.forEach(item => item.classList.remove('is-active'));
}
// htmx.on("htmx:responseError", function (evt) {
//     const requestConfig = evt.detail.requestConfig;
//     const error = evt.detail.xhr.response;
//     const msgElem = document.querySelectorAll("#notification-msg");
//     msgElem.forEach(item => item.innerHTML = error);
//     const notElem = document.querySelectorAll("#notification-bar");
//     notElem.forEach(item => item.classList.remove('is-hidden'));
// })
function hideNotification(elementId: string) {
    const element = document.querySelectorAll("#" + elementId);
    element.forEach(item => item.classList.add('is-hidden'));
}
export function toggleClass(root: ParentNode, elementId: string, cls = "is-active") {
    function toggleItem(item: Element) {
        if (item.classList.contains(cls)) {
            item.classList.remove(cls);
        } else {
            item.classList.add(cls);
        }
    }
    const element = root.querySelectorAll("#" + elementId);
    element.forEach(item => toggleItem(item));
}
