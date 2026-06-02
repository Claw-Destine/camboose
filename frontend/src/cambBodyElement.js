import htmx from 'htmx.org';
import { ShadowTemplElement } from './elementBase.js';
import bodyTemplate from './templates/camb-body.html'
import { registerElementWithTemplate } from "./elementBase.js";
import { toggleClass } from './utils.js';

export class CambBody extends ShadowTemplElement {
    static observedAttributes = ["project-id", "view"];
    projectId = this.getAttribute("project-id")
    view = this.getAttribute("view")

    constructor() {
        super("camb-body", true)
        const pid = this.getAttribute("project-id")
        const view = this.getAttribute("view")

        this.setupMenu(pid, view)
        this.setupQuickProjectMenu(pid, view)
    }

    setupMenu(currentPid, currentView) {
        const menuItems = this.shadowRoot.querySelectorAll('#menu-list a');

        function setActiveMenu(clickedElement) {
            menuItems.forEach(item => item.classList.remove('is-active'));
            clickedElement.classList.add('is-active');
        }

        function setupMenuItem(item) {
            const view = item.id.split("-")[0]
            item.addEventListener('click', event => setActiveMenu(event.target, view))
            let path = "/components/" + view
            if (currentPid) {
                path = path + "?currentProject=" + currentPid
            }
            item.setAttribute("hx-get", path)
            if (view == currentView) {
                setActiveMenu(item)
            }
        }

        menuItems.forEach(item => setupMenuItem(item));
    }

    setupQuickProjectMenu(currentPid, currentView) {
        let button = this.shadowRoot.querySelector('#quick-project-button');
        button.addEventListener('click', _ => toggleClass(this.shadowRoot, "quick-project-select"))
        let caption = button.querySelector("#quick-project-caption")

        function setupProjectLink(item, root) {
            item.addEventListener('click', _ => toggleClass(root, "quick-project-select"))
            let path = "/components/" + currentView
            if (currentPid) {
                path = path + "?currentProject=" + currentPid
            }
            item.setAttribute("hx-get", path)
            caption.textContent = item.textContent
        }

        // let menu = this.shadowRoot.querySelectorAll('#quick-project-menu')[0];
        let slot = this.shadowRoot.querySelector('#quick-project-slot')
        let projectLinks = slot.assignedNodes();
        projectLinks.forEach(item => setupProjectLink(item, this.shadowRoot));
    }

    attributeChangedCallback(name, oldValue, newValue) {
        console.log(`Attribute ${name} has changed.`);
    }
}

export function registerCambBody() {
    registerElementWithTemplate("camb-body", CambBody, bodyTemplate)
}
