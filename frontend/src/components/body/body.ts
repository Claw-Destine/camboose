import { ShadowTemplElement } from '../../elementBase';
import bodyTemplate from './camb-body.html';
import { registerElementWithTemplate } from "../../elementBase";
import { toggleClass } from '../../utils';

export class CambBody extends ShadowTemplElement {
    static observedAttributes = ["project-id", "view"];
    projectId = this.getAttribute("project-id");
    view = this.getAttribute("view");

    constructor() {
        super("camb-body", true);
        const pid = this.getAttribute("project-id");
        const view = this.getAttribute("view");

        this.setupMenu(pid, view);
        this.setupQuickProjectMenu();
    }

    setupMenu(currentPid: string | null, currentView: string | null) {
        const menuItems = this.shadowRoot?.querySelectorAll<HTMLAnchorElement>('#menu-list a');
        if (!menuItems) {
            return;
        }

        function setActiveMenu(clickedElement: Element) {
            menuItems.forEach(item => item.classList.remove('is-active'));
            clickedElement.classList.add('is-active');
        }

        function setupMenuItem(item: HTMLAnchorElement) {
            const view = item.id.split("-")[0];
            item.addEventListener('click', event => {
                const target = event.currentTarget;
                if (target instanceof Element) {
                    setActiveMenu(target);
                }
            });

            let path = "/components/" + view;
            if (currentPid) {
                path = path + "?currentProject=" + currentPid;
            }
            item.setAttribute("hx-get", path);
            if (view == currentView) {
                setActiveMenu(item);
            }
        }

        menuItems.forEach(item => setupMenuItem(item));
    }

    setupQuickProjectMenu() {
        const button = this.shadowRoot?.querySelector<HTMLElement>('#quick-project-button');
        if (!button || !this.shadowRoot) {
            return;
        }

        button.addEventListener('click', _ => toggleClass(this.shadowRoot as ShadowRoot, "quick-project-select"));
        const caption = button.querySelector<HTMLElement>("#quick-project-caption");
        if (!caption) {
            return;
        }

        function setCaption(item: Node, root: ShadowRoot) {
            if (item instanceof Element && item.classList.contains("is-active")) {
                caption.textContent = item.textContent;
            }

        }

        const slot = this.shadowRoot.querySelector<HTMLSlotElement>('#quick-project-slot');
        if (!slot) {
            return;
        }

        const projectLinks = slot.assignedNodes();
        projectLinks.forEach(item => setCaption(item, this.shadowRoot));
    }

    attributeChangedCallback(name: string, oldValue: string | null, newValue: string | null) {
        console.log(`Attribute ${name} has changed.`);
    }
}

export function registerCambBody() {
    registerElementWithTemplate("camb-body", CambBody, bodyTemplate);
}
