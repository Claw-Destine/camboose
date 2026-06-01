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

        const menuItems = this.querySelectorAll('#menu-list a');
        menuItems.forEach(item => item.addEventListener('click', event => this.setActiveMenu(event)));
    }

    setActiveMenu(clickedElement) {
        const menuItems = this.querySelectorAll('#menu-list a');
        menuItems.forEach(item => item.classList.remove('is-active'));
        clickedElement.classList.add('is-active');
    }
}

export function registerCambBody() {
    registerElementWithTemplate("camb-body", CambBody, bodyTemplate)
}
