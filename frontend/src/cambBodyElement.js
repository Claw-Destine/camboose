import htmx from 'htmx.org';
import { ShadowTemplElement } from './elementBase.js';
import bodyTemplate from './templates/camb-body.html'
import { registerElementWithTemplate } from "./elementBase.js";

export class CambBody extends ShadowTemplElement {
    static observedAttributes = ["project-id", "view"];

    constructor() {
        super("camb-body", true)
        const pid = this.getAttribute("project-id")
        if (!pid) {
            console.error("project-id not set for camb-body")
        }

        this.setupMenu()
    }

    setupMenu() {
        const menuItems = this.shadowRoot.querySelectorAll('#menu-list a');

        function setActiveMenu(clickedElement) {
            menuItems.forEach(item => item.classList.remove('is-active'));
            clickedElement.classList.add('is-active');
        }

        function setupMenuItem(item) {
            item.addEventListener('click', event => setActiveMenu(event.target))
            const path = "/components/" + item.id.split("-")[0]
            item.setAttribute("hx-get", path)
        }

        menuItems.forEach(item => setupMenuItem(item));
    }
}

export function registerCambBody() {
    registerElementWithTemplate("camb-body", CambBody, bodyTemplate)
}
