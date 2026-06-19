export function toggleClass(root: ParentNode, elementId: string, cls = "is-active") {
    function toggleItem(item: Element) {
        if (item.classList.contains(cls)) {
            item.classList.remove(cls);
        } else {
            item.classList.add(cls);
        }
    }
    const element = root.querySelectorAll("#" + elementId);
    element.forEach((item) => toggleItem(item));
}
export function insertCustomElement(innerHTML: string, root: Element) {
    root.innerHTML += innerHTML;
}
export function removeElement(query: string, root: ParentNode) {
    const el = root.querySelector(query) as Element;
    el.remove();
}
export function getCookie(name: string): string | null {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop()?.split(";").shift() || null;
    return null;
}
export function showNotification(event: CustomEvent) {
    const nbar = document.body.querySelector("#notification-bar");
    const msg = nbar.querySelector("#notification-msg");

    const status = event.detail.xhr.status;
    const text = event.detail.xhr.responseText;

    msg.textContent = "Request returned status: " + status + " . Reason: " + text;
    nbar.classList.remove("is-hidden");
}
export function hideNotification() {
    const nbar = document.body.querySelector("#notification-bar");

    nbar.classList.add("is-hidden");
}
