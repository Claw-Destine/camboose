import "bulma/css/bulma.min.css";
import htmx from 'htmx.org';
import { TemplElement } from './elementBase.js';
import bodyTemplate from './templates/camb-body.html'

export class CambBody extends TemplElement {
    static observedAttributes = ["project-id", "view"];

    constructor() {
        super(bodyTemplate)
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