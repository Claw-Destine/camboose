import { ShadowTemplElement } from '../../elementBase';
import projectsTemplate from './camb-projects.html';
import projectTemplate from './camb-project.html';
import { registerElementWithTemplate } from '../../elementBase';
import htmx from 'htmx.org';

export class ProjectsComponent extends ShadowTemplElement {
    constructor() {
        super("camb-projects", true);
        this.setupProjectLinks();
    }

    setupProjectLinks() {
        const slot = this.shadowRoot?.querySelector<HTMLSlotElement>('#projects-list');
        if (!slot) {
            return;
        }

        slot.addEventListener('click', event => {
            const path = event.composedPath();
            const link = path.find(node => node instanceof HTMLAnchorElement);
            if (!(link instanceof HTMLAnchorElement)) {
                return;
            }

            const url = link.getAttribute('data-project-url') || link.getAttribute('href');
            if (!url || url === '#') {
                return;
            }

            const target = this.shadowRoot?.querySelector<HTMLElement>('#project-details');
            if (!target) {
                return;
            }

            event.preventDefault();
            event.stopPropagation();
            htmx.ajax('GET', url, target);
        });
    }
}

export function registerProjectsComponent() {
    registerElementWithTemplate("camb-projects", ProjectsComponent, projectsTemplate);
}

export class ProjectComponent extends ShadowTemplElement {
    constructor() {
        super("camb-projects", true);
    }
}

export function registerProjectComponent() {
    registerElementWithTemplate("camb-project", ProjectComponent, projectTemplate);
}