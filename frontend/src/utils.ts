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
export function removeElement(query: string, root: Element) {
    const el = root.querySelector(query) as Element;
    el.remove();
}
export function getCookie(name: string): string | null {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop()?.split(";").shift() || null;
    return null;
}
