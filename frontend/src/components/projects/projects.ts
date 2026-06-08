import { ShadowTemplElement } from '../../elementBase';
import projectsTemplate from './camb-projects.html';
import projectTemplate from './camb-project.html';
import { registerElementWithTemplate } from '../../elementBase';
import htmx from 'htmx.org';

export class ProjectsComponent extends ShadowTemplElement {
    constructor() {
        super("camb-projects", true);
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